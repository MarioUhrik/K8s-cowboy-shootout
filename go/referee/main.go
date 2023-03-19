package main

import (
	pb "github.com/MarioUhrik/K8s-cowboy-shootout/go/proto/pb"
	"context"
	"log"
	"time"
	"os"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Cowboy struct {
	name string `json:"name"`
	health int32 `json:"health"`
	damage int32 `json:"damage"`
}

var k8sConfig *rest.Config
var k8sClientset *kubernetes.Clientset
var cowboys []Cowboy
var namespace string
var cowboysReady bool = false
var winnerFound bool = false

func initK8s() {
	log.Printf("Initializing")
	k8sConfig, _ = rest.InClusterConfig()
	k8sClientset, _ = kubernetes.NewForConfig(k8sConfig)

	namespace = os.Getenv("K8S_NAMESPACE") // TODO - define env variable in k8s manifest
	log.Printf("Initialized")
}

func listPods() *v1.PodList {
	listOptions := meta_v1.ListOptions{LabelSelector: "microservice=cowboy"}
	podList, err := k8sClientset.CoreV1().Pods(namespace).List(context.TODO(), listOptions)
	if err != nil {
		panic(err)
	}
	return podList
}

func getRemainingCowboyIPs() ([]string) {
	var podIPs []string
	for _, pod := range listPods().Items {
		if (pod.Status.ContainerStatuses[0].Ready) {
			podIPs = append(podIPs, pod.Status.PodIP)
		}
	}
	return podIPs
}

func waitForReadiness() {
	log.Printf("Waiting for cowboy readiness")
	for (!cowboysReady) {
		if (len(getRemainingCowboyIPs()) == len(cowboys)) {
			cowboysReady = true
			break
		}
		time.Sleep(1 * time.Second)
	}
	log.Printf("Cowboys are ready")
}

func startDuel() {
	for _, cowboyIP := range getRemainingCowboyIPs() { // TODO: first establish all connections, then call RPCs all at the same time
		cowboyURL := cowboyIP + ":8080"
		conn, err := grpc.Dial(cowboyURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)
		}
		client := pb.NewCowboyClient(conn)
		log.Printf("Ordering cowboy %s to start shooting", cowboyURL)
		client.StartShooting(context.Background(), &pb.Empty{})
		conn.Close()
	}
}

func findWinner() {
	for (!winnerFound) {
		cowboyIPs := getRemainingCowboyIPs();
		if (len(cowboyIPs) == 1) {
			cowboyURL := cowboyIPs[0] + ":8080"
			conn, err := grpc.Dial(cowboyURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				panic(err)
			}
			client := pb.NewCowboyClient(conn)
			log.Printf("Declaring cowboy %s the winner", cowboyURL)
			client.GetDeclaredVictorious(context.Background(), &pb.Empty{})
			conn.Close()
			winnerFound = true
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	initK8s()
	waitForReadiness()
	startDuel()
	findWinner()
	// TODO - helm templating for cowboy input file
	// TODO - readiness probe for cowboys
}
