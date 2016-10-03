package main

import (
	
	"fmt"
	"time"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

type AWSResponse struct {
	Datapoints []struct {
		Average     float64     `json:"Average"`
		Maximum     interface{} `json:"Maximum"`
		Minimum     interface{} `json:"Minimum"`
		SampleCount interface{} `json:"SampleCount"`
		Sum         interface{} `json:"Sum"`
		Timestamp   string      `json:"Timestamp"`
		Unit        string      `json:"Unit"`
		
	} `json:"Datapoints"`
	Label string `json:"Label"`
	InstanceId  string      `json:"InstanceId"`
}

func main() {
    
	    r := httprouter.New()
		r.GET("/metrics", getMetrics)
		r.GET("/cpu" ,getCPUUtilization)
		r.GET("/netwrkinput",getNetworkInput)
		r.GET("/netwrkoutput",getNetworkOutput)
		r.GET("/diskreadops",getDiskReadOps)
		r.GET("/diskreadbytes",getDiskReadBytes)
		r.GET("/diskwritebytes",getDiskWriteBytes)
		r.GET("/memory", getMemoryUtilization)
		r.GET("/count", getHTTPCount)
		http.ListenAndServe("localhost:8080",r)

}

func getMetrics(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			
   
	svc := ec2.New(session.New())

		resp, err := svc.DescribeInstances(nil)
		if err != nil {
        panic(err)
    }

    // resp has all of the response data, pull out instance IDs:
	fmt.Println("The Response data is", resp)
    fmt.Println("> Number of reservation sets: ", len(resp.Reservations))
    for idx, res := range resp.Reservations {
			fmt.Println("  > Number of instances: ", len(res.Instances))
			
			for _, inst := range resp.Reservations[idx].Instances {
				fmt.Println("    - Instance ID: ", *inst.InstanceId)
					

			}
     }

	  jsonResp, _ := json.Marshal(resp)
	  w.Header().Set("Content-Type", "application/json")
	  w.Header().Set("Access-Control-Allow-Origin", "*")
	  w.WriteHeader(200)
	  fmt.Fprintf(w, "%s", jsonResp)
	 

}

func getCPUUtilization(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	 svc := ec2.New(session.New())

		resp, err := svc.DescribeInstances(nil)
		if err != nil {
        panic(err)
    }
	
	var sliceArray []AWSResponse

	
	
    // resp has all of the response data, pull out instance IDs:
	//fmt.Println("The Response data is", resp)
    fmt.Println("> Number of reservation sets: ", len(resp.Reservations))
    for idx, res := range resp.Reservations {
			fmt.Println("  > Number of instances: ", len(res.Instances))
			
			for _, inst := range resp.Reservations[idx].Instances {
				fmt.Println("    - Instance ID: ", *inst.InstanceId)
				
				cw := cloudwatch.New(session.New())

					params := &cloudwatch.GetMetricStatisticsInput{
					EndTime:    aws.Time(time.Now()),     // Required
					MetricName: aws.String("CPUUtilization"), // Required
					Namespace:  aws.String("AWS/EC2"),  // Required
					Period:     aws.Int64(3600),             // Required
					StartTime:  aws.Time(time.Now().Add(-120 *time.Minute)),     // Required
					Statistics: []*string{ // Required
						aws.String("Average"), // Required
						// More values...
					},
					Dimensions: []*cloudwatch.Dimension{
						{ // Required
							Name:  aws.String("InstanceId"),  // Required
							Value: aws.String(*inst.InstanceId), // Required
						},
						// More values...
					},
					Unit: aws.String("Percent"),
				  }
				  
				  fmt.Println("Params are", params)
				
				  metrics, err := cw.GetMetricStatistics(params)

				   if err != nil {
					fmt.Println("Error")
					return
				 }

				 fmt.Println(metrics)
				 
				  AWS := AWSResponse{}
				  jsonResp, _ := json.Marshal(metrics)
				  err1 := json.Unmarshal(jsonResp, &AWS)
				  fmt.Println("Instance id is", *inst.InstanceId)
				  AWS.InstanceId = *inst.InstanceId
				  fmt.Println("Instance ID inside is" , AWS.InstanceId)
				  sliceArray = append(sliceArray, AWS)
				  if err1!= nil {
					  panic(err1)
					} 
 
				 
			}
     }
	  fmt.Println("Final array is", sliceArray)
	  finalResponse, _ :=json.Marshal(sliceArray)
	  w.Header().Set("Content-Type", "application/json")
	  w.Header().Set("Access-Control-Allow-Origin", "*")
	  w.WriteHeader(200)
	  fmt.Fprintf(w, "%s", finalResponse)
}	


func getNetworkInput(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	 svc := ec2.New(session.New())

		resp, err := svc.DescribeInstances(nil)
		if err != nil {
        panic(err)
    }
	
	var sliceArray []AWSResponse

	
	
    // resp has all of the response data, pull out instance IDs:
	//fmt.Println("The Response data is", resp)
    fmt.Println("> Number of reservation sets: ", len(resp.Reservations))
    for idx, res := range resp.Reservations {
			fmt.Println("  > Number of instances: ", len(res.Instances))
			
			for _, inst := range resp.Reservations[idx].Instances {
				fmt.Println("    - Instance ID: ", *inst.InstanceId)
				
				cw := cloudwatch.New(session.New())
				
				networkIn := &cloudwatch.GetMetricStatisticsInput{
				EndTime:    aws.Time(time.Now()),     // Required
				MetricName: aws.String("NetworkIn"), // Required
				Namespace:  aws.String("AWS/EC2"),  // Required
				Period:     aws.Int64(3600),             // Required
				StartTime:  aws.Time(time.Now().Add(-120 * time.Minute)),     // Required
				Statistics: []*string{ // Required
					aws.String("Average"), // Required
					// More values...
				},
				Dimensions: []*cloudwatch.Dimension{
					{ // Required
						Name:  aws.String("InstanceId"),  // Required
						Value: aws.String(*inst.InstanceId), // Required
					},
					// More values...
				},
				Unit: aws.String("Bytes"),
			  }
			  
				  networkInMetrics, err := cw.GetMetricStatistics(networkIn)

				   if err != nil {
					fmt.Println("Error")
					return
				 }

				 fmt.Println(networkInMetrics)
			
				 
				  AWS := AWSResponse{}
				  jsonResp, _ := json.Marshal(networkInMetrics)
				  err1 := json.Unmarshal(jsonResp, &AWS)
				  fmt.Println("Instance id is", *inst.InstanceId)
				  AWS.InstanceId = *inst.InstanceId
				  fmt.Println("Instance ID inside is" , AWS.InstanceId)
				  sliceArray = append(sliceArray, AWS)
				  if err1!= nil {
					  panic(err1)
					} 
 
				 
			}
     }
	  fmt.Println("Final array is", sliceArray)
	  finalResponse, _ :=json.Marshal(sliceArray)
	  w.Header().Set("Content-Type", "application/json")
	  w.Header().Set("Access-Control-Allow-Origin", "*")
	  w.WriteHeader(200)
	  fmt.Fprintf(w, "%s", finalResponse)
}				

func getNetworkOutput(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	 svc := ec2.New(session.New())

		resp, err := svc.DescribeInstances(nil)
		if err != nil {
        panic(err)
    }
	
	var sliceArray []AWSResponse

	
	
    // resp has all of the response data, pull out instance IDs:
	//fmt.Println("The Response data is", resp)
    fmt.Println("> Number of reservation sets: ", len(resp.Reservations))
    for idx, res := range resp.Reservations {
			fmt.Println("  > Number of instances: ", len(res.Instances))
			
			for _, inst := range resp.Reservations[idx].Instances {
				fmt.Println("    - Instance ID: ", *inst.InstanceId)
				
				cw := cloudwatch.New(session.New())
				
				networkOut := &cloudwatch.GetMetricStatisticsInput{
				EndTime:    aws.Time(time.Now()),     // Required
				MetricName: aws.String("NetworkOut"), // Required
				Namespace:  aws.String("AWS/EC2"),  // Required
				Period:     aws.Int64(3600),             // Required
				StartTime:  aws.Time(time.Now().Add(-120 * time.Minute)),     // Required
				Statistics: []*string{ // Required
					aws.String("Average"), // Required
					// More values...
				},
				Dimensions: []*cloudwatch.Dimension{
					{ // Required
						Name:  aws.String("InstanceId"),  // Required
						Value: aws.String(*inst.InstanceId), // Required
					},
					// More values...
				},
				Unit: aws.String("Bytes"),
			  }
			  
			  networkOutMetrics, err := cw.GetMetricStatistics(networkOut)

			   if err != nil {
				fmt.Println("Error")
				return
			 }

			 fmt.Println(networkOutMetrics)
			
				 
				  AWS := AWSResponse{}
				  jsonResp, _ := json.Marshal(networkOutMetrics)
				  err1 := json.Unmarshal(jsonResp, &AWS)
				  fmt.Println("Instance id is", *inst.InstanceId)
				  AWS.InstanceId = *inst.InstanceId
				  fmt.Println("Instance ID inside is" , AWS.InstanceId)
				  sliceArray = append(sliceArray, AWS)
				  if err1!= nil {
					  panic(err1)
					} 
 
				 
			}
     }
	  fmt.Println("Final array is", sliceArray)
	  finalResponse, _ :=json.Marshal(sliceArray)
	  w.Header().Set("Content-Type", "application/json")
	  w.Header().Set("Access-Control-Allow-Origin", "*")
	  w.WriteHeader(200)
	  fmt.Fprintf(w, "%s", finalResponse)
}				
	
func getDiskReadOps(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	 svc := ec2.New(session.New())

		resp, err := svc.DescribeInstances(nil)
		if err != nil {
        panic(err)
    }
	
	var sliceArray []AWSResponse

	
	
    // resp has all of the response data, pull out instance IDs:
	//fmt.Println("The Response data is", resp)
    fmt.Println("> Number of reservation sets: ", len(resp.Reservations))
    for idx, res := range resp.Reservations {
			fmt.Println("  > Number of instances: ", len(res.Instances))
			
			for _, inst := range resp.Reservations[idx].Instances {
				fmt.Println("    - Instance ID: ", *inst.InstanceId)
				
				cw := cloudwatch.New(session.New())
				
				
				diskReadOps := &cloudwatch.GetMetricStatisticsInput{
				EndTime:    aws.Time(time.Now()),     // Required
				MetricName: aws.String("DiskReadOps"), // Required
				Namespace:  aws.String("AWS/EC2"),  // Required
				Period:     aws.Int64(3600),             // Required
				StartTime:  aws.Time(time.Now().Add(-120 * time.Minute)),     // Required
				Statistics: []*string{ // Required
					aws.String("Average"), // Required
					// More values...
				},
				Dimensions: []*cloudwatch.Dimension{
					{ // Required
						Name:  aws.String("InstanceId"),  // Required
						Value: aws.String(*inst.InstanceId), // Required
					},
					// More values...
				},
				Unit: aws.String("Count"),
			  }
			  
			  diskReadOpsMetrics, err := cw.GetMetricStatistics(diskReadOps)

			   if err != nil {
				fmt.Println("Error")
				return
			 }

			 fmt.Println(diskReadOpsMetrics)
			 
			
				 
				  AWS := AWSResponse{}
				  jsonResp, _ := json.Marshal(diskReadOpsMetrics)
				  err1 := json.Unmarshal(jsonResp, &AWS)
				  fmt.Println("Instance id is", *inst.InstanceId)
				  AWS.InstanceId = *inst.InstanceId
				  fmt.Println("Instance ID inside is" , AWS.InstanceId)
				  sliceArray = append(sliceArray, AWS)
				  if err1!= nil {
					  panic(err1)
					} 
 
				 
			}
     }
	  fmt.Println("Final array is", sliceArray)
	  finalResponse, _ :=json.Marshal(sliceArray)
	  w.Header().Set("Content-Type", "application/json")
	  w.Header().Set("Access-Control-Allow-Origin", "*")
	  w.WriteHeader(200)
	  fmt.Fprintf(w, "%s", finalResponse)
}				
	
func getDiskReadBytes(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	 svc := ec2.New(session.New())

		resp, err := svc.DescribeInstances(nil)
		if err != nil {
        panic(err)
    }
	
	var sliceArray []AWSResponse

	
	
    // resp has all of the response data, pull out instance IDs:
	//fmt.Println("The Response data is", resp)
    fmt.Println("> Number of reservation sets: ", len(resp.Reservations))
    for idx, res := range resp.Reservations {
			fmt.Println("  > Number of instances: ", len(res.Instances))
			
			for _, inst := range resp.Reservations[idx].Instances {
				fmt.Println("    - Instance ID: ", *inst.InstanceId)
				
				cw := cloudwatch.New(session.New())
				
				diskReadBytes := &cloudwatch.GetMetricStatisticsInput{
		 
				EndTime:    aws.Time(time.Now()),     // Required
				MetricName: aws.String("DiskReadBytes"), // Required
				Namespace:  aws.String("AWS/EC2"),  // Required
				Period:     aws.Int64(3600),             // Required
				StartTime:  aws.Time(time.Now().Add(-120 * time.Minute)),     // Required
				Statistics: []*string{ // Required
					aws.String("Average"), // Required
					// More values...
				},
				Dimensions: []*cloudwatch.Dimension{
					{ // Required
						Name:  aws.String("InstanceId"),  // Required
						Value: aws.String(*inst.InstanceId), // Required
					},
					// More values...
				},
				Unit: aws.String("Bytes"),
			  }
			  
			  diskReadBytesMetrics, err := cw.GetMetricStatistics(diskReadBytes)

			   if err != nil {
				fmt.Println("Error")
				return
			 }

			 fmt.Println(diskReadBytesMetrics)

			
				 
				  AWS := AWSResponse{}
				  jsonResp, _ := json.Marshal(diskReadBytesMetrics)
				  err1 := json.Unmarshal(jsonResp, &AWS)
				  fmt.Println("Instance id is", *inst.InstanceId)
				  AWS.InstanceId = *inst.InstanceId
				  fmt.Println("Instance ID inside is" , AWS.InstanceId)
				  sliceArray = append(sliceArray, AWS)
				  if err1!= nil {
					  panic(err1)
					} 
 
				 
			}
     }
	  fmt.Println("Final array is", sliceArray)
	  finalResponse, _ :=json.Marshal(sliceArray)
	  w.Header().Set("Content-Type", "application/json")
	  w.Header().Set("Access-Control-Allow-Origin", "*")
	  w.WriteHeader(200)
	  fmt.Fprintf(w, "%s", finalResponse)
}				
	
func getDiskWriteBytes(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	 svc := ec2.New(session.New())

		resp, err := svc.DescribeInstances(nil)
		if err != nil {
        panic(err)
    }
	
	var sliceArray []AWSResponse

	
	
    // resp has all of the response data, pull out instance IDs:
	//fmt.Println("The Response data is", resp)
    fmt.Println("> Number of reservation sets: ", len(resp.Reservations))
    for idx, res := range resp.Reservations {
			fmt.Println("  > Number of instances: ", len(res.Instances))
			
			for _, inst := range resp.Reservations[idx].Instances {
				fmt.Println("    - Instance ID: ", *inst.InstanceId)
				
				cw := cloudwatch.New(session.New())
				diskWriteBytes := &cloudwatch.GetMetricStatisticsInput{
					EndTime:    aws.Time(time.Now()),     // Required
					MetricName: aws.String("DiskWriteBytes"), // Required
					Namespace:  aws.String("AWS/EC2"),  // Required
					Period:     aws.Int64(3600),             // Required
					StartTime:  aws.Time(time.Now().Add(-120 * time.Minute)),     // Required
					Statistics: []*string{ // Required
						aws.String("Average"), // Required
						// More values...
					},
					Dimensions: []*cloudwatch.Dimension{
						{ // Required
							Name:  aws.String("InstanceId"),  // Required
							Value: aws.String(*inst.InstanceId), // Required
						},
						// More values...
					},
					Unit: aws.String("Bytes"),
				  }		  
				  
					
				  diskWriteBytesMetrics, err := cw.GetMetricStatistics(diskWriteBytes)

				   if err != nil {
					fmt.Println("Error")
					return
				 }

				 fmt.Println(diskWriteBytesMetrics)
			
				 
				  AWS := AWSResponse{}
				  jsonResp, _ := json.Marshal(diskWriteBytesMetrics)
				  err1 := json.Unmarshal(jsonResp, &AWS)
				  fmt.Println("Instance id is", *inst.InstanceId)
				  AWS.InstanceId = *inst.InstanceId
				  fmt.Println("Instance ID inside is" , AWS.InstanceId)
				  sliceArray = append(sliceArray, AWS)
				  if err1!= nil {
					  panic(err1)
					} 
 
				 
			}
     }
	  fmt.Println("Final array is", sliceArray)
	  finalResponse, _ :=json.Marshal(sliceArray)
	  w.Header().Set("Content-Type", "application/json")
	  w.Header().Set("Access-Control-Allow-Origin", "*")
	  w.WriteHeader(200)
	  fmt.Fprintf(w, "%s", finalResponse)
}

	
func getMemoryUtilization(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	 svc := ec2.New(session.New())

		resp, err := svc.DescribeInstances(nil)
		if err != nil {
        panic(err)
    }
	
	var sliceArray []AWSResponse

	
	
    // resp has all of the response data, pull out instance IDs:
	//fmt.Println("The Response data is", resp)
    fmt.Println("> Number of reservation sets: ", len(resp.Reservations))
    for idx, res := range resp.Reservations {
			fmt.Println("  > Number of instances: ", len(res.Instances))
			
			for _, inst := range resp.Reservations[idx].Instances {
				fmt.Println("    - Instance ID: ", *inst.InstanceId)
				
				cw := cloudwatch.New(session.New())
				memoryStats := &cloudwatch.GetMetricStatisticsInput{
					EndTime:    aws.Time(time.Now()),     // Required
					MetricName: aws.String("MemoryUtilization"), // Required
					Namespace:  aws.String("System/Linux"),  // Required
					Period:     aws.Int64(3600),             // Required
					StartTime:  aws.Time(time.Now().Add(-120 * time.Minute)),     // Required
					Statistics: []*string{ // Required
						aws.String("Average"), // Required
						// More values...
					},
					Dimensions: []*cloudwatch.Dimension{
						{ // Required
							Name:  aws.String("InstanceId"),  // Required
							Value: aws.String(*inst.InstanceId), // Required
						},
						// More values...
					},
					
				  }		  
				  
				  fmt.Println("Memory Stats are",memoryStats)
					
				  memoryStatsResponse, err := cw.GetMetricStatistics(memoryStats)

				   if err != nil {
					fmt.Println("Inside Error")
					return
				 }

				 fmt.Println(memoryStatsResponse)
			
				 
				  AWS := AWSResponse{}
				  jsonResp, _ := json.Marshal(memoryStatsResponse)
				  err1 := json.Unmarshal(jsonResp, &AWS)
				  fmt.Println("Instance id is", *inst.InstanceId)
				  AWS.InstanceId = *inst.InstanceId
				  fmt.Println("Instance ID inside is" , AWS.InstanceId)
				  sliceArray = append(sliceArray, AWS)
				  if err1!= nil {
					  panic(err1)
					} 
 
				 
			}
     }
	  fmt.Println("Final array is", sliceArray)
	  finalResponse, _ :=json.Marshal(sliceArray)
	  w.Header().Set("Content-Type", "application/json")
	  w.Header().Set("Access-Control-Allow-Origin", "*")
	  w.WriteHeader(200)
	  fmt.Fprintf(w, "%s", finalResponse)
}

func getHTTPCount(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			
   
	 svc := ec2.New(session.New())

		resp, err := svc.DescribeInstances(nil)
		if err != nil {
        panic(err)
    }
	
	var sliceArray []AWSResponse

	
	
    // resp has all of the response data, pull out instance IDs:
	//fmt.Println("The Response data is", resp)
    fmt.Println("> Number of reservation sets: ", len(resp.Reservations))
    for idx, res := range resp.Reservations {
			fmt.Println("  > Number of instances: ", len(res.Instances))
			
			for _, inst := range resp.Reservations[idx].Instances {
				fmt.Println("    - Instance ID: ", *inst.InstanceId)
				
				cw := cloudwatch.New(session.New())
				fmt.Println("CW is", cw)
				memoryStats := &cloudwatch.GetMetricStatisticsInput{
					EndTime:    aws.Time(time.Now()),     // Required
					MetricName: aws.String("RequestCount"), // Required
					Namespace:  aws.String("AWS/ELB"),  // Required
					Period:     aws.Int64(3600),             // Required
					StartTime:  aws.Time(time.Now().Add(-120 * time.Minute)),     // Required
					Statistics: []*string{ // Required
						aws.String("Sum"), // Required
						// More values...
					},
					Dimensions: []*cloudwatch.Dimension{
						{ // Required
							Name:  aws.String("LoadBalancerName"),  // Required
							Value: aws.String("PlsWork"), // Required
						},
						// More values...
					},
					
				  }		  
				  
				 // fmt.Println("Memory Stats are",memoryStats)
					
				  memoryStatsResponse, err := cw.GetMetricStatistics(memoryStats)

				   if err != nil {
					fmt.Println("Inside Error")
					return
				 }

				 fmt.Println(memoryStatsResponse)
			
				 
				  AWS := AWSResponse{}
				  jsonResp, _ := json.Marshal(memoryStatsResponse)
				  err1 := json.Unmarshal(jsonResp, &AWS)
				  fmt.Println("Instance id is", *inst.InstanceId)
				  AWS.InstanceId = *inst.InstanceId
				  fmt.Println("Instance ID inside is" , AWS.InstanceId)
				  sliceArray = append(sliceArray, AWS)
				  if err1!= nil {
					  panic(err1)
					} 
 
				 
			}
     }
	  fmt.Println("Final array is", sliceArray)
	  finalResponse, _ :=json.Marshal(sliceArray)
	  w.Header().Set("Content-Type", "application/json")
	  w.Header().Set("Access-Control-Allow-Origin", "*")
	  w.WriteHeader(200)
	  fmt.Fprintf(w, "%s", finalResponse)
	 

}			
		