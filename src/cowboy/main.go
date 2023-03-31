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
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type server struct {
	pb.UnimplementedCowboyServer
}

type Cowboy struct {
	name   string
	health int32
	damage int32
}

var cowboy Cowboy
var s grpc.Server
var healthServer *health.Server
var triggerShutdown chan string

var k8sConfig *rest.Config
var k8sClientset *kubernetes.Clientset
var namespace string

func (s *server) GetShot(ctx context.Context, request *pb.GetShotRequest) (*pb.GetShotResponse, error) {
	if cowboy.name == request.ShooterName {
		log.Printf("%s didn't hit anyone", cowboy.name)
		return &pb.GetShotResponse{VictimName: cowboy.name, RemainingHealth: cowboy.health}, nil
	}

	log.Printf("%s got shot by %s", cowboy.name, request.ShooterName)
	cowboy.health = cowboy.health - request.IncomingDamage
	log.Printf("%s has %d health left", cowboy.name, cowboy.health)
	if cowboy.health <= 0 {
		die()
	}

	return &pb.GetShotResponse{VictimName: cowboy.name, RemainingHealth: cowboy.health}, nil
}

func listPods() *v1.PodList {
	listOptions := meta_v1.ListOptions{LabelSelector: "microservice=cowboy"}
	log.Printf("DEBUG: About to request the list of pods")
	podList, err := k8sClientset.CoreV1().Pods(namespace).List(context.TODO(), listOptions)
	if err != nil {
		log.Panicf("Failed to list the cowboy pods: %v", err)
	}
	return podList
}

func getRemainingCowboyIPs() []string {
	var podIPs []string
	for _, pod := range listPods().Items {
		if pod.Status.ContainerStatuses[0].Ready {
			podIPs = append(podIPs, pod.Status.PodIP)
		}
	}
	return podIPs
}

func isVictorious() bool {
	return cowboy.health > 0 && len(getRemainingCowboyIPs()) == 1
}

func die() {
	log.Printf("%s is dead", cowboy.name)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
	time.Sleep(1 * time.Second) // avoid io timeout error on connections from other cowboys at this time
	triggerShutdown <- "Shutting down the GRPC server"
}

func shoot() {
	conn, err := grpc.Dial("cowboys:8080", grpc.WithTimeout(100*time.Millisecond), grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to Dial cowboy while shooting: %v", err)
		return
	}
	client := pb.NewCowboyClient(conn)
	log.Printf("%s shoots", cowboy.name)
	_, err = client.GetShot(context.Background(), &pb.GetShotRequest{ShooterName: cowboy.name, IncomingDamage: cowboy.damage})
	if err != nil {
		log.Printf("Failed to hit target cowboy while shooting: %v", err)
	}
	conn.Close()
}

func getReady() {
	k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Panicf("Failed to load the InClusterConfig for Kubernetes: %v", err)
	}
	k8sConfig.Timeout = time.Second
	k8sClientset, err = kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		log.Panicf("Failed to use the InClusterConfig for Kubernetes: %v", err)
	}
	namespace = os.Getenv("K8S_NAMESPACE")

	cowboyHealth, err := strconv.Atoi(os.Getenv("COWBOY_HEALTH"))
	if err != nil {
		log.Panicf("Failed to parse COWBOY_HEALTH env variable: %v", err)
	}
	cowboyDamage, err := strconv.Atoi(os.Getenv("COWBOY_DAMAGE"))
	if err != nil {
		log.Panicf("Failed to parse COWBOY_DAMAGE env variable: %v", err)
	}

	cowboy.name = os.Getenv("COWBOY_NAME")
	cowboy.health = int32(cowboyHealth)
	cowboy.damage = int32(cowboyDamage)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Panicf("Failed to listen on TCP port: %v", err)
	}

	triggerShutdown = make(chan string)
	s := grpc.NewServer()
	go func() {
		pb.RegisterCowboyServer(s, &server{})
		healthServer = health.NewServer()
		healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
		healthgrpc.RegisterHealthServer(s, healthServer)
		if err := s.Serve(listener); err != nil {
			log.Panicf("Failed to serve: %v", err)
		}
		<-triggerShutdown
		s.GracefulStop()
		listener.Close()
	}()
	log.Printf("%s is now ready", cowboy.name)
}

func waitForReadiness() {
	ready := false
	for !ready {
		time.Sleep(1 * time.Second)
		if len(getRemainingCowboyIPs()) == len(listPods().Items) {
			ready = true
		}
	}
}

func shootout() {
	for cowboy.health > 0 {
		log.Printf("DEBUG: About to shoot")
		shoot()
		log.Printf("DEBUG: About to check for victory")
		time.Sleep(1000 * time.Millisecond)
		if isVictorious() {
			log.Printf("%s is victorious! The fastest hand in the West.", cowboy.name)
			return
		}
	}
}

func main() {
	getReady()
	waitForReadiness()
	shootout()
	time.Sleep(3600 * time.Second)
}
