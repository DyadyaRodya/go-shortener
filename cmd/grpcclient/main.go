package main

import (
	"context"
	"log"
	"net/url"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/DyadyaRodya/go-shortener/proto/v1"
)

func main() {
	ctx := context.Background()
	conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewGoShortenerServiceClient(conn)
	res, err := c.Ping(ctx, &pb.PingRequest{})
	log.Println(status.FromError(err))
	log.Println(res)

	res1, err := c.BatchCreateShortURL(ctx, &pb.BatchCreateShortURLRequest{
		Urls: []*pb.BatchCreateShortURLRequestItem{
			&pb.BatchCreateShortURLRequestItem{
				CorrelationId: "1",
				OriginalUrl:   "https://google.com",
			},
			&pb.BatchCreateShortURLRequestItem{
				CorrelationId: "2",
				OriginalUrl:   "https://ya.ru",
			},
			&pb.BatchCreateShortURLRequestItem{
				CorrelationId: "3",
				OriginalUrl:   "https://mail.ru",
			},
		},
	})
	log.Println(status.FromError(err))
	log.Println(res1)

	md := metadata.New(map[string]string{"auth": res1.GetNewJwtToken()})
	ctx = metadata.NewOutgoingContext(ctx, md)

	res2, err := c.GetUserShortURLs(ctx, &pb.GetUserShortURLsRequest{})
	log.Println(status.FromError(err))
	log.Println(res2)

	md = metadata.New(map[string]string{"auth": res2.GetNewJwtToken()})
	ctx = metadata.NewOutgoingContext(ctx, md)

	res3, err := c.GetStats(ctx, &pb.GetStatsRequest{RealIp: "192.168.1.1"})
	log.Println(status.FromError(err))
	log.Println(res3)

	res3, err = c.GetStats(ctx, &pb.GetStatsRequest{RealIp: "192.168.10.1"})
	log.Println(status.FromError(err)) // Permission Denied
	log.Println(res3)

	res4, err := c.CreateShortURL(ctx, &pb.CreateShortURLRequest{Url: "https://ya3.ru"})
	log.Println(status.FromError(err))
	log.Println(res4)

	parsed, _ := url.Parse(res4.Result)
	id := strings.TrimPrefix(parsed.Path, "/")

	res4, err = c.CreateShortURL(ctx, &pb.CreateShortURLRequest{Url: "https://ya3.ru"})
	log.Println(status.FromError(err))
	log.Println(res4) // already exists

	res5, err := c.GetFullByID(ctx, &pb.GetFullByIDRequest{Id: id})
	log.Println(status.FromError(err))
	log.Println(res5)

	res6, err := c.DeleteShortURLs(ctx, &pb.DeleteShortURLsRequest{Ids: []string{id}})
	log.Println(status.FromError(err))
	log.Println(res6)

	time.Sleep(11 * time.Second)

	res5, err = c.GetFullByID(ctx, &pb.GetFullByIDRequest{Id: id})
	log.Println(status.FromError(err))
	log.Println(res5) // Gone
}
