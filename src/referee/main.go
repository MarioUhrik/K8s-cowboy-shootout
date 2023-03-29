package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/MarioUhrik/K8s-cowboy-shootout/src/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var k8sConfig *rest.Config
var k8sClientset *kubernetes.Clientset
var namespace string
var cowboysReady bool = false
var winnerFound bool = false

func initK8s() {
	log.Printf("Initializing")
	k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Panicf("Failed to load the InClusterConfig for Kubernetes: %v", err)
	}
	k8sClientset, err = kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		log.Panicf("Failed to use the InClusterConfig for Kubernetes: %v", err)
	}
	namespace = os.Getenv("K8S_NAMESPACE")
	log.Printf("Initialized")
}

func listPods() *v1.PodList {
	listOptions := meta_v1.ListOptions{LabelSelector: "microservice=cowboy"}
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

func waitForReadiness() {
	log.Printf("Waiting for cowboy readiness")
	for !cowboysReady {
		if len(getRemainingCowboyIPs()) == len(listPods().Items) {
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
		log.Printf("Ordering cowboy %s to start shooting", cowboyURL)
		conn, err := grpc.Dial(cowboyURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("Failed to Dial cowboy %s: %v", cowboyURL, err)
			continue
		}
		client := pb.NewCowboyClient(conn)
		_, err = client.StartShooting(context.Background(), &pb.StartShootingRequest{})
		if err != nil {
			log.Printf("Failed to order cowboy %s to start shooting: %v", cowboyURL, err)
			continue
		}
		err = conn.Close()
		if err != nil {
			log.Printf("Failed to close connection to cowboy %s after ordering him to start shooting: %v", cowboyURL, err)
			continue
		}
	}
}

func findWinner() {
	for !winnerFound {
		cowboyIPs := getRemainingCowboyIPs()
		if len(cowboyIPs) == 1 {
			cowboyURL := cowboyIPs[0] + ":8080"
			conn, err := grpc.Dial(cowboyURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Panicf("Failed to Dial cowboy %s: %v", cowboyURL, err)
			}
			client := pb.NewCowboyClient(conn)
			log.Printf("Declaring cowboy %s the winner", cowboyURL)
			_, err = client.GetDeclaredVictorious(context.Background(), &pb.GetDeclaredVictoriousRequest{})
			if err != nil {
				log.Panicf("Failed to declare cowboy %s victorious: %v", cowboyURL, err)
			}
			conn.Close()
			if err != nil {
				log.Printf("Failed to close connection to cowboy %s after declaring him victorious: %v", cowboyURL, err)
			}
			winnerFound = true
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func postShootout() { // TODO: use Kubernetes jobs or a retrypolicy=never instead?
	for true {
		time.Sleep(3600 * time.Second)
	}
}

func main() { // TODO: Theoretically, all functionalities of the referee could be done by the cowboys themselves instead. Decomission the Referee
	initK8s()
	waitForReadiness()
	startDuel()
	findWinner()
	postShootout()
}
