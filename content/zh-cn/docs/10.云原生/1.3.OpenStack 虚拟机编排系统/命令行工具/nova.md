---
title: nova
---


语法格式：nova  \[OPTIONS]  \[SubCommand  \[OPTIONS]]

nova list \[OPTIONS] # 列出SERVER相关信息

nova show  # 显示指定SERVER的详细信息，非常详细

nova instacne-action-list  # 列出指定SERVER的操作，创建、启动、停止、删除时间等，

语法格式：nova  \[OPTIONS]  \[SubCommand  \[OPTIONS]]

注意：语法中的SERVER指的都是已经创建的虚拟服务器，SERVER可以用实例的NAME(实例名)或者UUID(实例的ID)来表示，SERVER的ID和NAME可以用过nova list命令查到

UUID:Universally Unique Identifier,即通用唯一识别码

可以使用nova help SubCommand命令查看相关子命令的使用方法

nova list \[OPTIONS] # 列出SERVER相关信息

1. OPTIONS

2. --all-tenants # 显示所有租户的SERVER信息，可简写为--all-t

3. --tenant  \[] # 显示指定租户的SERVER信息

4. EXAMPLE

5. nova list --all-t # 显示所有正在运行的实例，可以查看实例以及ID和主机名

6. nova list --all-t --host `cat /etc/uuid` # 显示`cat /etc/uuid`命令输出的主机名的节点运行的实例信息

nova show  # 显示指定SERVER的详细信息，非常详细

1. EXAMPLE

2. nova show ID # 以实例ID展示该实例的详细信息

3. nova show ID | grep host # 以实例ID查看所在节点的主机名

nova instacne-action-list  # 列出指定SERVER的操作，创建、启动、停止、删除时间等，

1. EXAMPLE

2. nova instance-action-list ID|NAME # 以实例ID显示该实例的活动信息，包括启动、停止、创建时间等

nova

nova \[--version] \[--debug] \[--os-cache] \[--timings]

            [--os-region-name ] [--service-type ]

            [--service-name ]

            [--os-endpoint-type ]

            [--os-compute-api-version ]

            [--endpoint-override ] [--profile HMAC_KEY]

            [--insecure] [--os-cacert ]

            [--os-cert ] [--os-key ] [--timeout ]

            [--os-auth-type ] [--os-auth-url OS_AUTH_URL]

            [--os-system-scope OS_SYSTEM_SCOPE] [--os-domain-id OS_DOMAIN_ID]

            [--os-domain-name OS_DOMAIN_NAME] [--os-project-id OS_PROJECT_ID]

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

Command-line interface to the OpenStack Nova API.

Positional arguments:位置参数

子命令

    add-secgroup                Add a Security Group to a server.

    agent-create                Create new agent build.

    agent-delete                Delete existing agent build.

    agent-list                  List all builds.

    agent-modify                Modify existing agent build.

    aggregate-add-host          Add the host to the specified aggregate.

    aggregate-create            Create a new aggregate with the specified

                                details.

    aggregate-delete            Delete the aggregate.

    aggregate-list              Print a list of all aggregates.

    aggregate-remove-host       Remove the specified host from the specified

                                aggregate.

    aggregate-set-metadata      Update the metadata associated with the

                                aggregate.

    aggregate-show              Show details of the specified aggregate.

    aggregate-update            Update the aggregate's name and optionally

                                availability zone.

    availability-zone-list      List all the availability zones.

    backup                      Backup a server by creating a 'backup' type

                                snapshot.

    boot                        Boot a new server.

    cell-capacities             Get cell capacities for all cells or a given

                                cell.

    cell-show                   Show details of a given cell.

    clear-password              Clear the admin password for a server from the

                                metadata server. This action does not actually

                                change the instance server password.

    console-log                 Get console log output of a server.

    delete                      Immediately shut down and delete specified

                                server(s).

    diagnostics                 Retrieve server diagnostics.

    evacuate                    Evacuate server from failed host.

    flavor-access-add           Add flavor access for the given tenant.

    flavor-access-list          Print access information about the given

                                flavor.

    flavor-access-remove        Remove flavor access for the given tenant.

    flavor-create               Create a new flavor.

    flavor-delete               Delete a specific flavor

    flavor-key                  Set or unset extra_spec for a flavor.

    flavor-list                 Print a list of available 'flavors' (sizes of

                                servers).

    flavor-show                 Show details about the given flavor.

    flavor-update               Update the description of an existing flavor.

                                (Supported by API versions '2.55' -

                                '2.latest') [hint: use '--os-compute-api-

                                version' flag to show help message for proper

                                version]

    force-delete                Force delete a server.

    get-mks-console             Get an MKS console to a server. (Supported by

                                API versions '2.8' - '2.latest') [hint: use

                                '--os-compute-api-version' flag to show help

                                message for proper version]

    get-password                Get the admin password for a server. This

                                operation calls the metadata service to query

                                metadata information and does not read

                                password information from the server itself.

    get-rdp-console             Get a rdp console to a server.

    get-serial-console          Get a serial console to a server.

    get-spice-console           Get a spice console to a server.

    get-vnc-console             Get a vnc console to a server.

    host-evacuate               Evacuate all instances from failed host.

    host-evacuate-live          Live migrate all instances of the specified

                                host to other available hosts.

    host-meta                   Set or Delete metadata on all instances of a

                                host.

    host-servers-migrate        Cold migrate all instances off the specified

                                host to other available hosts.

    hypervisor-list             List hypervisors. (Supported by API versions

                                '2.0' - '2.latest') [hint: use '--os-compute-

                                api-version' flag to show help message for

                                proper version]

    hypervisor-servers          List servers belonging to specific

                                hypervisors.

    hypervisor-show             Display the details of the specified

                                hypervisor.

    hypervisor-stats            Get hypervisor statistics over all compute

                                nodes.

    hypervisor-uptime           Display the uptime of the specified

                                hypervisor.

    image-create                Create a new image by taking a snapshot of a

                                running server.

    instance-action             Show an action.

    instance-action-list        List actions on a server. (Supported by API

                                versions '2.0' - '2.latest') [hint: use '--os-

                                compute-api-version' flag to show help message

                                for proper version]

    interface-attach            Attach a network interface to a server.

    interface-detach            Detach a network interface from a server.

    interface-list              List interfaces attached to a server.

    keypair-add                 Create a new key pair for use with servers.

    keypair-delete              Delete keypair given by its name. (Supported

                                by API versions '2.0' - '2.latest') [hint: use

                                '--os-compute-api-version' flag to show help

                                message for proper version]

    keypair-list                Print a list of keypairs for a user (Supported

                                by API versions '2.0' - '2.latest') [hint: use

                                '--os-compute-api-version' flag to show help

                                message for proper version]

    keypair-show                Show details about the given keypair.

                                (Supported by API versions '2.0' - '2.latest')

                                [hint: use '--os-compute-api-version' flag to

                                show help message for proper version]

    limits                      Print rate and absolute limits.

    list                        List servers.

    list-extensions             List all the os-api extensions that are

                                available.

    list-secgroup               List Security Group(s) of a server.

    live-migration              Migrate running server to a new machine.

    live-migration-abort        Abort an on-going live migration. (Supported

                                by API versions '2.24' - '2.latest') [hint:

                                use '--os-compute-api-version' flag to show

                                help message for proper version]

    live-migration-force-complete

                                Force on-going live migration to complete.

                                (Supported by API versions '2.22' -

                                '2.latest') [hint: use '--os-compute-api-

                                version' flag to show help message for proper

                                version]

    lock                        Lock a server. A normal (non-admin) user will

                                not be able to execute actions on a locked

                                server.

    meta                        Set or delete metadata on a server.

    migrate                     Migrate a server.

    migration-list              Print a list of migrations. (Supported by API

                                versions '2.0' - '2.latest') [hint: use '--os-

                                compute-api-version' flag to show help message

                                for proper version]

    pause                       Pause a server.

    quota-class-show            List the quotas for a quota class.

    quota-class-update          Update the quotas for a quota class.

                                (Supported by API versions '2.0' - '2.latest')

                                [hint: use '--os-compute-api-version' flag to

                                show help message for proper version]

    quota-defaults              List the default quotas for a tenant.

    quota-delete                Delete quota for a tenant/user so their quota

                                will Revert back to default.

    quota-show                  List the quotas for a tenant/user.

    quota-update                Update the quotas for a tenant/user.

                                (Supported by API versions '2.0' - '2.latest')

                                [hint: use '--os-compute-api-version' flag to

                                show help message for proper version]

    reboot                      Reboot a server.

    rebuild                     Shutdown, re-image, and re-boot a server.

    refresh-network             Refresh server network information.

    remove-secgroup             Remove a Security Group from a server.

    rescue                      Reboots a server into rescue mode, which

                                starts the machine from either the initial

                                image or a specified image, attaching the

                                current boot disk as secondary.

    reset-network               Reset network of a server.

    reset-state                 Reset the state of a server.

    resize                      Resize a server.

    resize-confirm              Confirm a previous resize.

    resize-revert               Revert a previous resize (and return to the

                                previous VM).

    restore                     Restore a soft-deleted server.

    resume                      Resume a server.

    server-group-create         Create a new server group with the specified

                                details.

    server-group-delete         Delete specific server group(s).

    server-group-get            Get a specific server group.

    server-group-list           Print a list of all server groups.

    server-migration-list       Get the migrations list of specified server.

                                (Supported by API versions '2.23' -

                                '2.latest') [hint: use '--os-compute-api-

                                version' flag to show help message for proper

                                version]

    server-migration-show       Get the migration of specified server.

                                (Supported by API versions '2.23' -

                                '2.latest') [hint: use '--os-compute-api-

                                version' flag to show help message for proper

                                version]

    server-tag-add              Add one or more tags to a server. (Supported

                                by API versions '2.26' - '2.latest') [hint:

                                use '--os-compute-api-version' flag to show

                                help message for proper version]

    server-tag-delete           Delete one or more tags from a server.

                                (Supported by API versions '2.26' -

                                '2.latest') [hint: use '--os-compute-api-

                                version' flag to show help message for proper

                                version]

    server-tag-delete-all       Delete all tags from a server. (Supported by

                                API versions '2.26' - '2.latest') [hint: use

                                '--os-compute-api-version' flag to show help

                                message for proper version]

    server-tag-list             Get list of tags from a server. (Supported by

                                API versions '2.26' - '2.latest') [hint: use

                                '--os-compute-api-version' flag to show help

                                message for proper version]

    server-tag-set              Set list of tags to a server. (Supported by

                                API versions '2.26' - '2.latest') [hint: use

                                '--os-compute-api-version' flag to show help

                                message for proper version]

    service-delete              Delete the service by UUID ID. (Supported by

                                API versions '2.0' - '2.latest') [hint: use

                                '--os-compute-api-version' flag to show help

                                message for proper version]

    service-disable             Disable the service. (Supported by API

                                versions '2.0' - '2.latest') [hint: use '--os-

                                compute-api-version' flag to show help message

                                for proper version]

    service-enable              Enable the service. (Supported by API versions

                                '2.0' - '2.latest') [hint: use '--os-compute-

                                api-version' flag to show help message for

                                proper version]

    service-force-down          Force service to down. (Supported by API

                                versions '2.11' - '2.latest') [hint: use

                                '--os-compute-api-version' flag to show help

                                message for proper version]

    service-list                Show a list of all running services. Filter by

                                host & binary.

    set-password                Change the admin password for a server.

    shelve                      Shelve a server.

    shelve-offload              Remove a shelved server from the compute node.

    show                        Show details about the given server.

    ssh                         SSH into a server.

    start                       Start the server(s).

    stop                        Stop the server(s).

    suspend                     Suspend a server.

    trigger-crash-dump          Trigger crash dump in an instance. (Supported by API versions '2.17' - '2.latest') [hint:

                                use '--os-compute-api-version' flag to show

                                help message for proper version]

    unlock                      Unlock a server.

    unpause                     Unpause a server.

    unrescue                    Restart the server from normal boot disk

                                again.

    unshelve                    Unshelve a server.

    update                      Update the name or the description for a

                                server.

    usage                       Show usage data for a single tenant.

    usage-list                  List usage data for all tenants.

    version-list                List all API versions.

    volume-attach               Attach a volume to a server.

    volume-attachments          List all the volumes attached to a server.

    volume-detach               Detach a volume from a server.

    volume-update               Update the attachment on the server. Migrates

                                the data from an attached volume to the

                                specified available volume and swaps out the

                                active attachment to the new volume.

    bash-completion             Prints all of the commands and options to

                                stdout so that the nova.bash_completion script

                                doesn't have to hard code them.

    help                        Display help about this program or one of its

                                subcommands.

Optional arguments:可选参数

--version                     show program's version number and exit

--debug                       Print debugging output.

--os-cache                    Use the auth token cache. Defaults to False if

                                env[OS_CACHE] is not set.

--timings                     Print call timing info.

--os-region-name

                                Defaults to env[OS_REGION_NAME].

--service-type

                                Defaults to compute for most actions.

--service-name

                                Defaults to env[NOVA_SERVICE_NAME].

--os-endpoint-type

                                Defaults to env[NOVA_ENDPOINT_TYPE],

                                env[OS_ENDPOINT_TYPE] or publicURL.

--os-compute-api-version

                                Accepts X, X.Y (where X is major and Y is

                                minor part) or "X.latest", defaults to

                                env[OS_COMPUTE_API_VERSION].

--endpoint-override

                                Use this API endpoint instead of the Service

                                Catalog. Defaults to

                                env[NOVACLIENT_ENDPOINT_OVERRIDE].

--profile HMAC\_KEY            HMAC key to use for encrypting context data

                                for performance profiling of operation. This

                                key should be the value of the HMAC key

                                configured for the OSprofiler middleware in

                                nova; it is specified in the Nova

                                configuration file at "/etc/nova/nova.conf".

                                Without the key, profiling will not be

                                triggered even if OSprofiler is enabled on the

                                server side.

--os-auth-type , --os-auth-plugin

                                Authentication type to use

API Connection Options:API连接选项

Options controlling the HTTP API Connections

--insecure                    Explicitly allow client to perform "insecure"

                                TLS (https) requests. The server's certificate

                                will not be verified against any certificate

                                authorities. This option should be used with

                                caution.

--os-cacert   Specify a CA bundle file to use in verifying a

                                TLS (https) server certificate. Defaults to

                                env[OS_CACERT].

--os-cert        Defaults to env\[OS\_CERT].

--os-key                 Defaults to env\[OS\_KEY].

--timeout            Set request timeout (in seconds).

Authentication Options:认证选项

Options specific to the password plugin.

--os-auth-url OS\_AUTH\_URL     Authentication URL

--os-system-scope OS\_SYSTEM\_SCOPE

                                Scope for system operations

--os-domain-id OS\_DOMAIN\_ID   Domain ID to scope to

--os-domain-name OS\_DOMAIN\_NAME

                                Domain name to scope to

--os-project-id OS\_PROJECT\_ID, --os-tenant-id OS\_PROJECT\_ID

                                Project ID to scope to

--os-project-name OS\_PROJECT\_NAME, --os-tenant-name OS\_PROJECT\_NAME

                                Project name to scope to

--os-project-domain-id OS\_PROJECT\_DOMAIN\_ID

                                Domain ID containing project

--os-project-domain-name OS\_PROJECT\_DOMAIN\_NAME

                                Domain name containing project

--os-trust-id OS\_TRUST\_ID     Trust ID

--os-default-domain-id OS\_DEFAULT\_DOMAIN\_ID

                                Optional domain ID to use with v3 and v2

                                parameters. It will be used for both the user

                                and project domain in v3 and ignored in v2

                                authentication.

--os-default-domain-name OS\_DEFAULT\_DOMAIN\_NAME

                                Optional domain name to use with v3 API and v2

                                parameters. It will be used for both the user

                                and project domain in v3 and ignored in v2

                                authentication.

--os-user-id OS\_USER\_ID       User id

--os-username OS\_USERNAME, --os-user-name OS\_USERNAME

                                Username

--os-user-domain-id OS\_USER\_DOMAIN\_ID

                                User's domain id

--os-user-domain-name OS\_USER\_DOMAIN\_NAME

                                User's domain name

--os-password OS\_PASSWORD     User's password

See "nova help COMMAND" for help on a specific command.
