#!/bin/bash

# Define the input string
#input_string="postgresql://username:password@shelter-postgresql-primary:5432/shelter"
input_string=$POSTGRES_URL

# Extract the relevant parts of the string
username_password=$(echo $input_string | awk -F'[@]' '{print $1}')
host_port=$(echo $input_string | awk -F'[@]' '{print $2}' | awk -F'[,/]' '{print $1}')
database_name=$(echo $input_string | awk -F'[@/]' '{print $NF}')

# Construct the new string
output_string="${username_password}@${host_port}/${database_name}"

# Print the output
echo "Original string: $input_string"
echo "Converted string: $output_string"

migrate -verbose -path=/migrations -database ${output_string} up
