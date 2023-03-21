package main

import (
	pb "github.com/MarioUhrik/K8s-cowboy-shootout/go/proto/pb"
	"context"
	"log"
	"net"
	"time"
	"os"
	"strconv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	pb.UnimplementedCowboyServer
}

type Cowboy struct {
	name string
	health int32
	damage int32
	isInCombat bool
	isVictorious bool
}

var cowboy Cowboy
var s grpc.Server

func (s *server) GetShot(ctx context.Context, shooter *pb.Shooter) (*pb.Shooter, error) {
	if (cowboy.name == shooter.Name) {
		log.Printf("%s didn't hit anyone", cowboy.name)
		return shooter, nil
	}

	log.Printf("%s got shot by %s", cowboy.name, shooter.Name)
	cowboy.health = cowboy.health - shooter.Damage
	log.Printf("%s has %d health left", cowboy.name, cowboy.health)
	if (cowboy.health <= 0) {
		die()
	}

	return shooter, nil
}

func (s *server) StartShooting(ctx context.Context, empty *pb.Empty) (*pb.Empty, error) {
	cowboy.isInCombat = true
	return empty, nil
}

func (s *server) GetDeclaredVictorious(ctx context.Context, empty *pb.Empty) (*pb.Empty, error) {
	cowboy.isVictorious = true
	cowboy.isInCombat = false
	return empty, nil
}

func die() {
	log.Printf("%s is dead", cowboy.name)
	s.GracefulStop()
}

func shoot() {
	conn, err := grpc.Dial("cowboys:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client := pb.NewCowboyClient(conn)
	log.Printf("%s shoots", cowboy.name)
	client.GetShot(context.Background(), &pb.Shooter{cowboy.name, cowboy.health, cowboy.damage})
	conn.Close()
}

func getReady() {
	health, _ := strconv.Atoi(os.Getenv("COWBOY_HEALTH"))
	damage, _ := strconv.Atoi(os.Getenv("COWBOY_DAMAGE"))

	cowboy.name = os.Getenv("COWBOY_NAME")
	cowboy.health = int32(health)
	cowboy.damage = int32(damage)
	cowboy.isInCombat = false
	cowboy.isVictorious =  false

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	go func() error {
		pb.RegisterCowboyServer(s, &server{})
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
		return nil
	}()
	log.Printf("%s is now ready", cowboy.name)
}

func waitToStartShootout() {
	for (!cowboy.isInCombat) {
		log.Printf("%s is eagerly awaiting the shootout to begin", cowboy.name)
		time.Sleep(150 * time.Millisecond)
	}
}

func shootout() {
	for cowboy.health > 0 {
		shoot()
		time.Sleep(1000 * time.Millisecond)
		if (cowboy.isVictorious) {
			log.Printf("%s is victorious! The fastest hand in the West.", cowboy.name)
			return
		}
	}
}

func postShootout() {
	for true {
		time.Sleep(3600 * time.Second)
	}
}

func main() {
	getReady()
	waitToStartShootout()
	shootout()
	postShootout()
}
