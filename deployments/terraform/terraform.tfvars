# 3. terraform.tfvars
# This file sets the actual values for the variables. Add the following,
# replacing the placeholder values with your actual OCI details:

tenancy_ocid     = "ocid1.tenancy.oc1..example"
user_ocid        = "ocid1.user.oc1..example"
fingerprint      = "11:22:33:44:55:66:77:88:99:00:aa:bb:cc:dd:ee:ff"
private_key_path = "~/.oci/oci_api_key.pem"
region           = "us-ashburn-1"
compartment_ocid = "ocid1.compartment.oc1..example"
arm_ubuntu_image_id = "ocid1.image.oc1.iad.example"
amd_ubuntu_image_id = "ocid1.image.oc1.iad.example"

# Note: Keep this file secure and do not commit it to version control,
# as it contains sensitive information.