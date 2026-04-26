resource "google_container_cluster" "primary" {
  name     = "my-gke-cluster"
  location = "us-central1"

  resource_labels = {
    env  = "test"
    team = "devops"
  }
}

resource "google_container_cluster" "untagged" {
  name     = "my-gke-cluster-untagged"
  location = "us-central1"
}
