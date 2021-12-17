# Terraform Provider for Google Identity Platform

This is a Terraform provider which is used to configure Google Identity Platform, available from the Google Marketplace.

## Requirements

* [Terraform](https://www.terraform.io/downloads.html) 0.14+
* [Go](https://golang.org/doc/install) 1.16.0 or higher

## Installing the provider

Enter the provider directory and run the following command:

```shell
make install
```

## Using the provider

See the [example](./examples/main.tf) directory for an example usage.

## Importing Config

The resource ID used to import the Config must conform to the following syntax

`projects/<project-number>/config`

```terraform
terraform import identity_platform_config.test_config projects/1234567890/config
```
