package database

import (
	"context"
	"log"
	"time"

	"github.com/vamshireddy02/job-portal/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)


var connectionString string = "mongodb://localhost:27017/jobs"


type DB struct {
	client *mongo.Client
}

func Connect() *DB {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	return &DB {
		client: client,
	}
}

func (db *DB) GetJob(id string) *model.JobListing {

	jobCollec := db.client.Database("job-portal").Collection("jobs")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}

	var jobListing model.JobListing

	err := jobCollec.FindOne(ctx, filter).Decode(&jobListing)
	if err != nil {
		log.Fatal(err)
	}


	return &jobListing 
}

func (db *DB) GetJobs() []*model.JobListing {
	var jobCollec = db.client.Database("job-portal").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()


	var jobListings []*model.JobListing
	cursor, err := jobCollec.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	if err = cursor.All(context.TODO(), &jobListings);
	err != nil {
		panic(nil)
	}

	return jobListings
}

func (db *DB) CreateJobListing(jobInfo model.CreateJobListingInput) *model.JobListing {

	jobCollec := db.client.Database("job-portal").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Insert the job listing into the collection
	insertResult, err := jobCollec.InsertOne(ctx, bson.M{
		"title":       jobInfo.Title,
		"description": jobInfo.Description,
		"url":         jobInfo.URL,
		"company":     jobInfo.Company,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Retrieve the inserted job ID and return the new job listing
	insertedID := insertResult.InsertedID.(primitive.ObjectID).Hex()
	returnJobListing := model.JobListing{
		ID:          insertedID,
		Title:       jobInfo.Title,
		Company:     jobInfo.Company,
		Description: jobInfo.Description,
		URL:         jobInfo.URL,
	}
	return &returnJobListing
}

func (db *DB) UpdateJobListing(jobId string, jobInfo model.UpdateJobListingInput) *model.JobListing {
	jobCollec := db.client.Database("job-portal").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create an update document with the fields that need to be updated
	updateJobInfo := bson.M{}
	if jobInfo.Title != nil {
		updateJobInfo["title"] = jobInfo.Title
	}
	if jobInfo.Description != nil {
		updateJobInfo["description"] = jobInfo.Description
	}
	if jobInfo.URL != nil {
		updateJobInfo["url"] = jobInfo.URL
	}

	// Convert string job ID to ObjectID
	_id, _ := primitive.ObjectIDFromHex(jobId)
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": updateJobInfo} // Apply updates

	// Find and update the job listing
	results := jobCollec.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var jobListing model.JobListing
	// Decode the updated job listing
	if err := results.Decode(&jobListing); err != nil {
		log.Fatal(err)
	}

	return &jobListing
}

func (db *DB) DeleteJobListing(jobId string) *model.DeleteJobResponse {
	jobCollec := db.client.Database("job-portal").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Convert the string job ID to ObjectID
	_id, _ := primitive.ObjectIDFromHex(jobId)
	filter := bson.M{"_id": _id} // Filter for deleting the job

	// Delete the job from the collection
	_, err := jobCollec.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}

	// Return response containing the deleted job ID
	return &model.DeleteJobResponse{DeletedJobID: jobId}
}