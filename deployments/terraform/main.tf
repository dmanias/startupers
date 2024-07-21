# Configure the OCI provider
provider "oci" {
  tenancy_ocid     = var.tenancy_ocid
  user_ocid        = var.user_ocid
  fingerprint      = var.fingerprint
  private_key_path = var.private_key_path
  region           = var.region
}

# Variables
variable "tenancy_ocid" {}
variable "user_ocid" {}
variable "fingerprint" {}
variable "private_key_path" {}
variable "region" {}
variable "compartment_ocid" {}

# Networking
resource "oci_core_vcn" "kubernetes_vcn" {
  cidr_block     = "10.0.0.0/16"
  compartment_id = var.compartment_ocid
  display_name   = "Kubernetes VCN"
}

resource "oci_core_subnet" "kubernetes_subnet" {
  cidr_block        = "10.0.1.0/24"
  compartment_id    = var.compartment_ocid
  vcn_id            = oci_core_vcn.kubernetes_vcn.id
  display_name      = "Kubernetes Subnet"
  security_list_ids = [oci_core_security_list.kubernetes_security_list.id]
}

# Security List
resource "oci_core_security_list" "kubernetes_security_list" {
  compartment_id = var.compartment_ocid
  vcn_id         = oci_core_vcn.kubernetes_vcn.id
  display_name   = "Kubernetes Security List"

  egress_security_rules {
    destination = "0.0.0.0/0"
    protocol    = "all"
  }

  ingress_security_rules {
    protocol = "6"  # TCP
    source   = "0.0.0.0/0"

    tcp_options {
      min = 22
      max = 22
    }
  }

  ingress_security_rules {
    protocol = "6"  # TCP
    source   = "0.0.0.0/0"

    tcp_options {
      min = 6443
      max = 6443
    }
  }
}

# Instances
resource "oci_core_instance" "control_plane_1" {
  availability_domain = data.oci_identity_availability_domain.ad.name
  compartment_id      = var.compartment_ocid
  shape               = "VM.Standard.A1.Flex"
  display_name        = "Control Plane 1"

  shape_config {
    ocpus         = 2
    memory_in_gbs = 12
  }

  source_details {
    source_type = "image"
    source_id   = var.arm_ubuntu_image_id
  }

  create_vnic_details {
    subnet_id        = oci_core_subnet.kubernetes_subnet.id
    assign_public_ip = true
  }
}

resource "oci_core_instance" "control_plane_2" {
  availability_domain = data.oci_identity_availability_domain.ad.name
  compartment_id      = var.compartment_ocid
  shape               = "VM.Standard.A1.Flex"
  display_name        = "Control Plane 2"

  shape_config {
    ocpus         = 2
    memory_in_gbs = 12
  }

  source_details {
    source_type = "image"
    source_id   = var.arm_ubuntu_image_id
  }

  create_vnic_details {
    subnet_id        = oci_core_subnet.kubernetes_subnet.id
    assign_public_ip = true
  }
}

resource "oci_core_instance" "worker_1" {
  availability_domain = data.oci_identity_availability_domain.ad.name
  compartment_id      = var.compartment_ocid
  shape               = "VM.Standard.E2.1.Micro"
  display_name        = "Worker 1"

  source_details {
    source_type = "image"
    source_id   = var.amd_ubuntu_image_id
  }

  create_vnic_details {
    subnet_id        = oci_core_subnet.kubernetes_subnet.id
    assign_public_ip = true
  }
}

resource "oci_core_instance" "worker_2" {
  availability_domain = data.oci_identity_availability_domain.ad.name
  compartment_id      = var.compartment_ocid
  shape               = "VM.Standard.E2.1.Micro"
  display_name        = "Worker 2"

  source_details {
    source_type = "image"
    source_id   = var.amd_ubuntu_image_id
  }

  create_vnic_details {
    subnet_id        = oci_core_subnet.kubernetes_subnet.id
    assign_public_ip = true
  }
}

# Data sources
data "oci_identity_availability_domain" "ad" {
  compartment_id = var.tenancy_ocid
  ad_number      = 1
}

# Output
output "control_plane_1_public_ip" {
  value = oci_core_instance.control_plane_1.public_ip
}

output "control_plane_2_public_ip" {
  value = oci_core_instance.control_plane_2.public_ip
}

output "worker_1_public_ip" {
  value = oci_core_instance.worker_1.public_ip
}

output "worker_2_public_ip" {
  value = oci_core_instance.worker_2.public_ip
}