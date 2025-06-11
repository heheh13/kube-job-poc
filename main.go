package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var defaultCmd = "sleep 10"

func main() {
	http.HandleFunc("/start-job", handleStartJob)
	http.HandleFunc("/job-status", handleJobStatus)
	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}

func handleStartJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cmd := r.URL.Query().Get("cmd")
	fmt.Println(cmd)
	if cmd == "" {
		cmd = defaultCmd
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not get in-cluster config: %v", err), http.StatusInternalServerError)
		return
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not create clientset: %v", err), http.StatusInternalServerError)
		return
	}

	jobName := fmt.Sprintf("agent-job-%d", metav1.Now().Unix())
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:    "agent",
							Image:   os.Getenv("AGENT_IMAGE"),
							Command: []string{"/bin/sh", "-c", cmd},
						},
					},
				},
			},
		},
	}

	jobClient := clientset.BatchV1().Jobs("default")
	_, err = jobClient.Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create job: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Job %s created\n", jobName)
}

func handleJobStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	jobID := strings.TrimSpace(r.URL.Query().Get("id"))
	if jobID == "" {
		http.Error(w, "Missing job id", http.StatusBadRequest)
		return
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not get in-cluster config: %v", err), http.StatusInternalServerError)
		return
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not create clientset: %v", err), http.StatusInternalServerError)
		return
	}

	job, err := clientset.BatchV1().Jobs("default").Get(context.TODO(), jobID, metav1.GetOptions{})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get job: %v", err), http.StatusInternalServerError)
		return
	}

	status := job.Status
	response := make(map[string]interface{})
	response["job"] = jobID
	response["active"] = status.Active
	response["succeeded"] = status.Succeeded
	response["failed"] = status.Failed

	if job.Status.CompletionTime != nil && job.Status.StartTime != nil {
		duration := job.Status.CompletionTime.Sub(job.Status.StartTime.Time)
		response["completed"] = true
		response["duration"] = duration.String()
	} else {
		response["completed"] = false
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
