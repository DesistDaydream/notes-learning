---
title: glance
---

#

usage: glance \[--version] \[-d] \[-v] \[--get-schema] \[-f]

              [--os-image-url OS_IMAGE_URL]

              [--os-image-api-version OS_IMAGE_API_VERSION]

              [--profile HMAC_KEY] [--key-file OS_KEY] [--ca-file OS_CACERT]

              [--cert-file OS_CERT] [--os-region-name OS_REGION_NAME]

              [--os-auth-token OS_AUTH_TOKEN]

              [--os-service-type OS_SERVICE_TYPE]

              [--os-endpoint-type OS_ENDPOINT_TYPE] [--insecure]

              [--os-cacert ] [--os-cert ]

              [--os-key ] [--timeout ] [--os-auth-type ]

              [--os-auth-url OS_AUTH_URL] [--os-system-scope OS_SYSTEM_SCOPE]

              [--os-domain-id OS_DOMAIN_ID] [--os-domain-name OS_DOMAIN_NAME]

              [--os-project-id OS_PROJECT_ID]

              [--os-project-name OS_PROJECT_NAME]

              [--os-project-domain-id OS_PROJECT_DOMAIN_ID]

              [--os-project-domain-name OS_PROJECT_DOMAIN_NAME]

              [--os-trust-id OS_TRUST_ID]

              [--os-default-domain-id OS_DEFAULT_DOMAIN_ID]

              [--os-default-domain-name OS_DEFAULT_DOMAIN_NAME]

              [--os-user-id OS_USER_ID] [--os-username OS_USERNAME]

              [--os-user-domain-id OS_USER_DOMAIN_ID]

              [--os-user-domain-name OS_USER_DOMAIN_NAME]

              [--os-password OS_PASSWORD]

               ...

Command-line interface to the OpenStack Images API.

Positional arguments:

 explain Describe a specific model.

    image-create        Create a new image.

    image-create-via-import

                        EXPERIMENTAL: Create a new image via image import.

    image-deactivate    Deactivate specified image.

    image-delete        Delete specified image.

    image-download      Download a specific image.

    image-import        Initiate the image import taskflow.

    image-list          List images you can access.

    image-reactivate    Reactivate specified image.

    image-show          Describe a specific image.

    image-stage         Upload data for a specific image to staging.

    image-tag-delete    Delete the tag associated with the given image.

    image-tag-update    Update an image with the given tag.

    image-update        Update an existing image.

    image-upload        Upload data for a specific image.

    import-info         Print import methods available from Glance.

    location-add        Add a location (and related metadata) to an image.

    location-delete     Remove locations (and related metadata) from an image.

    location-update     Update metadata of an image's location.

    md-namespace-create

                        Create a new metadata definitions namespace.

    md-namespace-delete

                        Delete specified metadata definitions namespace with

                        its contents.

    md-namespace-import

                        Import a metadata definitions namespace from file or

                        standard input.

    md-namespace-list   List metadata definitions namespaces.

    md-namespace-objects-delete

                        Delete all metadata definitions objects inside a

                        specific namespace.

    md-namespace-properties-delete

                        Delete all metadata definitions property inside a

                        specific namespace.

    md-namespace-resource-type-list

                        List resource types associated to specific namespace.

    md-namespace-show   Describe a specific metadata definitions namespace.

    md-namespace-tags-delete

                        Delete all metadata definitions tags inside a specific

                        namespace.

    md-namespace-update

                        Update an existing metadata definitions namespace.

    md-object-create    Create a new metadata definitions object inside a

                        namespace.

    md-object-delete    Delete a specific metadata definitions object inside a

                        namespace.

    md-object-list      List metadata definitions objects inside a specific

                        namespace.

    md-object-property-show

                        Describe a specific metadata definitions property

                        inside an object.

    md-object-show      Describe a specific metadata definitions object inside

                        a namespace.

    md-object-update    Update metadata definitions object inside a namespace.

    md-property-create  Create a new metadata definitions property inside a

                        namespace.

    md-property-delete  Delete a specific metadata definitions property inside

                        a namespace.

    md-property-list    List metadata definitions properties inside a specific

                        namespace.

    md-property-show    Describe a specific metadata definitions property

                        inside a namespace.

    md-property-update  Update metadata definitions property inside a

                        namespace.

    md-resource-type-associate

                        Associate resource type with a metadata definitions

                        namespace.

    md-resource-type-deassociate

                        Deassociate resource type with a metadata definitions

                        namespace.

    md-resource-type-list

                        List available resource type names.

    md-tag-create       Add a new metadata definitions tag inside a namespace.

    md-tag-create-multiple

                        Create new metadata definitions tags inside a

                        namespace.

    md-tag-delete       Delete a specific metadata definitions tag inside a

                        namespace.

    md-tag-list         List metadata definitions tags inside a specific

                        namespace.

    md-tag-show         Describe a specific metadata definitions tag inside a

                        namespace.

    md-tag-update       Rename a metadata definitions tag inside a namespace.

    member-create       Create member for a given image.

    member-delete       Delete image member.

    member-list         Describe sharing permissions by image.

    member-update       Update the status of a member for a given image.

    task-create         Create a new task.

    task-list           List tasks you can access.

    task-show           Describe a specific task.

    bash-completion     Prints arguments for bash_completion.

    help                Display help about this program or one of its

                        subcommands.

Optional arguments:

--version show program's version number and exit

-d, --debug Defaults to env\[GLANCECLIENT_DEBUG].

-v, --verbose Print more verbose output.

--get-schema Ignores cached copy and forces retrieval of schema

                        that generates portions of the help text. Ignored with

                        API version 1.

-f, --force Prevent select actions from requesting user

                        confirmation.

--os-image-url OS_IMAGE_URL

                        Defaults to env[OS_IMAGE_URL]. If the provided image

                        url contains a version number and `--os-image-api-

                        version` is omitted the version of the URL will be

                        picked as the image api version to use.

--os-image-api-version OS_IMAGE_API_VERSION

                        Defaults to env[OS_IMAGE_API_VERSION] or 2.

--profile HMAC_KEY HMAC key to use for encrypting context data for

                        performance profiling of operation. This key should be

                        the value of HMAC key configured in osprofiler

                        middleware in glance, it is specified in glance

                        configuration file at /etc/glance/glance-api.conf and

                        /etc/glance/glance-registry.conf. Without key the

                        profiling will not be triggered even if osprofiler is

                        enabled on server side. Defaults to env[OS_PROFILE].

--key-file OS_KEY DEPRECATED! Use --os-key.

--ca-file OS_CACERT DEPRECATED! Use --os-cacert.

--cert-file OS_CERT DEPRECATED! Use --os-cert.

--os-region-name OS_REGION_NAME

                        Defaults to env[OS_REGION_NAME].

--os-auth-token OS_AUTH_TOKEN

                        Defaults to env[OS_AUTH_TOKEN].

--os-service-type OS_SERVICE_TYPE

                        Defaults to env[OS_SERVICE_TYPE].

--os-endpoint-type OS_ENDPOINT_TYPE

                        Defaults to env[OS_ENDPOINT_TYPE].

--os-auth-type , --os-auth-plugin

                        Authentication type to use

API Connection Options:

Options controlling the HTTP API Connections

--insecure Explicitly allow client to perform "insecure" TLS

                        (https) requests. The server's certificate will not be

                        verified against any certificate authorities. This

                        option should be used with caution.

--os-cacert

                        Specify a CA bundle file to use in verifying a TLS

                        (https) server certificate. Defaults to

                        env[OS_CACERT].

--os-cert

                        Defaults to env[OS_CERT].

--os-key Defaults to env\[OS_KEY].

--timeout Set request timeout (in seconds).

Authentication Options:

Options specific to the password plugin.

--os-auth-url OS_AUTH_URL

                        Authentication URL

--os-system-scope OS_SYSTEM_SCOPE

                        Scope for system operations

--os-domain-id OS_DOMAIN_ID

                        Domain ID to scope to

--os-domain-name OS_DOMAIN_NAME

                        Domain name to scope to

--os-project-id OS_PROJECT_ID, --os-tenant-id OS_PROJECT_ID

                        Project ID to scope to

--os-project-name OS_PROJECT_NAME, --os-tenant-name OS_PROJECT_NAME

                        Project name to scope to

--os-project-domain-id OS_PROJECT_DOMAIN_ID

                        Domain ID containing project

--os-project-domain-name OS_PROJECT_DOMAIN_NAME

                        Domain name containing project

--os-trust-id OS_TRUST_ID

                        Trust ID

--os-default-domain-id OS_DEFAULT_DOMAIN_ID

                        Optional domain ID to use with v3 and v2 parameters.

                        It will be used for both the user and project domain

                        in v3 and ignored in v2 authentication.

--os-default-domain-name OS_DEFAULT_DOMAIN_NAME

                        Optional domain name to use with v3 API and v2

                        parameters. It will be used for both the user and

                        project domain in v3 and ignored in v2 authentication.

--os-user-id OS_USER_ID

                        User id

--os-username OS_USERNAME, --os-user-name OS_USERNAME

                        Username

--os-user-domain-id OS_USER_DOMAIN_ID

                        User's domain id

--os-user-domain-name OS_USER_DOMAIN_NAME

                        User's domain name

--os-password OS_PASSWORD

                        User's password

See "glance help COMMAND" for help on a specific command.

Run `glance --os-image-api-version 1 help` for v1 help
