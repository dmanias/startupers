# 2. variables.tf
# This file declares all the variables used in main.tf. Add the following:

variable "tenancy_ocid" {
  description = "The OCID of your tenancy"
  type        = string
}

variable "user_ocid" {
  description = "The OCID of the user"
  type        = string
}

variable "fingerprint" {
  description = "The fingerprint of the key"
  type        = string
}

variable "private_key_path" {
  description = "The path to the private key"
  type        = string
}

variable "region" {
  description = "The OCI region"
  type        = string
}

variable "compartment_ocid" {
  description = "The OCID of the compartment"
  type        = string
}

variable "arm_ubuntu_image_id" {
  description = "The OCID of the ARM Ubuntu image"
  type        = string
}

variable "amd_ubuntu_image_id" {
  description = "The OCID of the AMD Ubuntu image"
  type        = string
}