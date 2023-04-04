package main

import (
	"context"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	pb "github.com/MarioUhrik/K8s-cowboy-shootout/src/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Cowboy struct {
	pb.UnimplementedCowboyServer

	name   string
	health int32
	damage int32

	cowboyServer    *grpc.Server
	healthServer    *health.Server
	triggerShutdown chan string

	k8sConfig    *rest.Config
	k8sClientset *kubernetes.Clientset
	namespace    string
}

func (self *Cowboy) GetShot(ctx context.Context, request *pb.GetShotRequest) (*pb.GetShotResponse, error) {
	if self.name == request.ShooterName {
		log.Printf("%s didn't hit anyone", self.name)
		return &pb.GetShotResponse{VictimName: self.name, RemainingHealth: self.health}, nil
	}

	log.Printf("%s got shot by %s", self.name, request.ShooterName)
	self.health = self.health - request.IncomingDamage
	log.Printf("%s has %d health left", self.name, self.health)
	if self.health <= 0 {
		self.die()
	}

	return &pb.GetShotResponse{VictimName: self.name, RemainingHealth: self.health}, nil
}

func (self *Cowboy) listPods() *v1.PodList {
	listOptions := meta_v1.ListOptions{LabelSelector: "microservice=cowboy"}
	podList, err := self.k8sClientset.CoreV1().Pods(self.namespace).List(context.TODO(), listOptions)
	if err != nil {
		log.Panicf("Failed to list the cowboy pods: %v", err)
	}
	return podList
}

func (self *Cowboy) getRemainingCowboyIPs() []string {
	var podIPs []string
	for _, pod := range self.listPods().Items {
		if pod.Status.ContainerStatuses[0].Ready {
			podIPs = append(podIPs, pod.Status.PodIP)
		}
	}
	return podIPs
}

func (self *Cowboy) isVictorious() bool {
	return self.health > 0 && len(self.getRemainingCowboyIPs()) == 1
}

func (self *Cowboy) die() {
	log.Printf("%s is dead", self.name)
	self.healthServer.Shutdown()
	time.Sleep(1 * time.Second) // avoid io timeout error on connections from other cowboys at this time
	self.triggerShutdown <- "Shutting down the GRPC server"
}

func (self *Cowboy) shoot() {
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
	conn, err := grpc.DialContext(ctx, "cowboys:8080", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("Failed to Dial cowboy while shooting: %v", err)
		return
	}
	client := pb.NewCowboyClient(conn)
	log.Printf("%s shoots", self.name)
	_, err = client.GetShot(ctx, &pb.GetShotRequest{ShooterName: self.name, IncomingDamage: self.damage})
	if err != nil {
		log.Printf("Failed to hit target cowboy while shooting: %v", err)
	}
	conn.Close()
}

func (self *Cowboy) getReady() {
	var err error
	self.k8sConfig, err = rest.InClusterConfig()
	if err != nil {
		log.Panicf("Failed to load the InClusterConfig for Kubernetes: %v", err)
	}
	self.k8sConfig.Timeout = time.Second
	self.k8sClientset, err = kubernetes.NewForConfig(self.k8sConfig)
	if err != nil {
		log.Panicf("Failed to use the InClusterConfig for Kubernetes: %v", err)
	}
	self.namespace = os.Getenv("K8S_NAMESPACE")

	cowboyHealth, err := strconv.Atoi(os.Getenv("COWBOY_HEALTH"))
	if err != nil {
		log.Panicf("Failed to parse COWBOY_HEALTH env variable: %v", err)
	}
	cowboyDamage, err := strconv.Atoi(os.Getenv("COWBOY_DAMAGE"))
	if err != nil {
		log.Panicf("Failed to parse COWBOY_DAMAGE env variable: %v", err)
	}

	self.name = os.Getenv("COWBOY_NAME")
	self.health = int32(cowboyHealth)
	self.damage = int32(cowboyDamage)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Panicf("Failed to listen on TCP port: %v", err)
	}

	self.triggerShutdown = make(chan string)
	self.cowboyServer = grpc.NewServer()
	go func() {
		pb.RegisterCowboyServer(self.cowboyServer, self)
		self.healthServer = health.NewServer()
		healthgrpc.RegisterHealthServer(self.cowboyServer, self.healthServer)
		if err := self.cowboyServer.Serve(listener); err != nil {
			log.Panicf("Failed to serve: %v", err)
		}

		<-self.triggerShutdown
		self.cowboyServer.GracefulStop()
		listener.Close()
	}()
	time.Sleep(5 * time.Second)
	log.Printf("%s is now ready", self.name)
}

func (self *Cowboy) waitForReadiness() {
	ready := false
	for !ready {
		time.Sleep(1 * time.Second)
		if len(self.getRemainingCowboyIPs()) == len(self.listPods().Items) {
			ready = true
		}
	}
}

func (self *Cowboy) shootout() {
	for self.health > 0 {
		self.shoot()
		time.Sleep(1000 * time.Millisecond)
		if self.isVictorious() {
			log.Printf("%s is victorious! The fastest hand in the West.", self.name)
			time.Sleep(3600 * time.Second)
		}
	}
}

func main() {
	cowboy := Cowboy{}
	cowboy.getReady()
	cowboy.waitForReadiness()
	cowboy.shootout()
}
