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
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	pb.UnimplementedCowboyServer
}

type Cowboy struct {
	name         string
	health       int32
	damage       int32
	isInCombat   bool
	isVictorious bool
}

var cowboy Cowboy
var s grpc.Server

func (s *server) GetShot(ctx context.Context, request *pb.GetShotRequest) (*pb.GetShotResponse, error) {
	if cowboy.name == request.ShooterName {
		log.Printf("%s didn't hit anyone", cowboy.name)
		return &pb.GetShotResponse{VictimName: cowboy.name, RemainingHealth: cowboy.health}, nil
	}

	log.Printf("%s got shot by %s", cowboy.name, request.ShooterName)
	cowboy.health = cowboy.health - request.IncomingDamage
	log.Printf("%s has %d health left", cowboy.name, cowboy.health)
	if cowboy.health <= 0 {
		die() // This may cause a segmentation fault, possibly because we're shutting down the server while it's answering requests
	} // This is potentially fixable by having die() set a global variable flag, and having the main function poll the value of that flag, calling GracefulStop() there

	return &pb.GetShotResponse{VictimName: cowboy.name, RemainingHealth: cowboy.health}, nil
}

func (s *server) StartShooting(ctx context.Context, request *pb.StartShootingRequest) (*pb.StartShootingResponse, error) {
	cowboy.isInCombat = true
	return &pb.StartShootingResponse{}, nil
}

func (s *server) GetDeclaredVictorious(ctx context.Context, request *pb.GetDeclaredVictoriousRequest) (*pb.GetDeclaredVictoriousResponse, error) {
	cowboy.isVictorious = true
	cowboy.isInCombat = false
	return &pb.GetDeclaredVictoriousResponse{}, nil
}

func die() {
	log.Printf("%s is dead", cowboy.name)
	s.Stop()
}

func shoot() {
	conn, err := grpc.Dial("cowboys:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to Dial cowboy while shooting: %v", err)
		return
	}
	client := pb.NewCowboyClient(conn)
	log.Printf("%s shoots", cowboy.name)
	_, err = client.GetShot(context.Background(), &pb.GetShotRequest{ShooterName: cowboy.name, IncomingDamage: cowboy.damage})
	if err != nil {
		log.Panicf("Failed to hit target cowboy while shooting: %v", err)
	}
	conn.Close()
}

func getReady() {
	health, err := strconv.Atoi(os.Getenv("COWBOY_HEALTH"))
	if err != nil {
		log.Panicf("Failed to parse COWBOY_HEALTH env variable: %v", err)
	}
	damage, err := strconv.Atoi(os.Getenv("COWBOY_DAMAGE"))
	if err != nil {
		log.Panicf("Failed to parse COWBOY_DAMAGE env variable: %v", err)
	}

	cowboy.name = os.Getenv("COWBOY_NAME")
	cowboy.health = int32(health)
	cowboy.damage = int32(damage)
	cowboy.isInCombat = false
	cowboy.isVictorious = false

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Panicf("Failed to listen on TCP port: %v", err)
	}

	s := grpc.NewServer()
	go func() {
		pb.RegisterCowboyServer(s, &server{})
		if err := s.Serve(listener); err != nil {
			log.Panicf("Failed to serve: %v", err)
		}
	}()
	log.Printf("%s is now ready", cowboy.name)
}

func waitToStartShootout() {
	for !cowboy.isInCombat {
		log.Printf("%s is eagerly awaiting the shootout to begin", cowboy.name)
		time.Sleep(150 * time.Millisecond)
	}
}

func shootout() {
	for cowboy.health > 0 {
		shoot()
		time.Sleep(1000 * time.Millisecond)
		if cowboy.isVictorious {
			log.Printf("%s is victorious! The fastest hand in the West.", cowboy.name)
			time.Sleep(3600 * time.Second)
			return
		}
	}
}

func main() {
	getReady()
	waitToStartShootout()
	shootout()
}
