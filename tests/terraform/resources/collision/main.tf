resource "aws_autoscaling_group" "aurora_cluster_bastion_auto_scaling_group" {
  default_cooldown          = "300"
  desired_capacity          = "1"
  force_delete              = "false"
  health_check_grace_period = "60"
  health_check_type         = "EC2"
  max_instance_lifetime     = "0"
  max_size                  = "1"
  metrics_granularity       = "1Minute"
  min_size                  = "1"
  name                      = "bc-aurora-cluster-bastion-auto-scaling-group"
  protect_from_scale_in     = "false"

  wait_for_capacity_timeout = "10m"


  tag {
    key                 = "Name"
    propagate_at_launch = true
    value               = "bc-aurora-bastion"
  }

  tag {
    key                 = "Env"
    propagate_at_launch = false
    value               = "prod"
  }

  tags = {
    git_org   = "bridgecrewio"
    git_repo  = "platform"
    yor_trace = "48564943-4cfc-403c-88cd-cbb207e0d33e"
  }
