package auth

// Code generated by Makefile; DO NOT EDIT.

var authModel = `{"schema_version":"1.1","type_definitions":[{"type":"user"},{"metadata":{"relations":{"member":{"directly_related_user_types":[{"type":"user"}]}}},"relations":{"member":{"this":{}}},"type":"group"},{"metadata":{"relations":{"can_edit":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_view":{},"server":{"directly_related_user_types":[{"type":"server"}]}}},"relations":{"can_edit":{"union":{"child":[{"this":{}},{"tupleToUserset":{"computedUserset":{"relation":"admin"},"tupleset":{"relation":"server"}}}]}},"can_view":{"tupleToUserset":{"computedUserset":{"relation":"viewer"},"tupleset":{"relation":"server"}}},"server":{"this":{}}},"type":"certificate"},{"metadata":{"relations":{"can_edit":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_view":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"project":{"directly_related_user_types":[{"type":"project"}]}}},"relations":{"can_edit":{"union":{"child":[{"this":{}},{"tupleToUserset":{"computedUserset":{"relation":"operator"},"tupleset":{"relation":"project"}}}]}},"can_view":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"can_edit"}},{"tupleToUserset":{"computedUserset":{"relation":"viewer"},"tupleset":{"relation":"project"}}}]}},"project":{"this":{}}},"type":"image"},{"metadata":{"relations":{"can_edit":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_view":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"project":{"directly_related_user_types":[{"type":"project"}]}}},"relations":{"can_edit":{"union":{"child":[{"this":{}},{"tupleToUserset":{"computedUserset":{"relation":"operator"},"tupleset":{"relation":"project"}}}]}},"can_view":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"can_edit"}},{"tupleToUserset":{"computedUserset":{"relation":"viewer"},"tupleset":{"relation":"project"}}}]}},"project":{"this":{}}},"type":"image_alias"},{"metadata":{"relations":{"admin":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_access_console":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_access_files":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_connect_sftp":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_edit":{},"can_exec":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_manage_backups":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_manage_snapshots":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_update_state":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_view":{},"operator":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"project":{"directly_related_user_types":[{"type":"project"}]},"user":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"viewer":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]}}},"relations":{"admin":{"union":{"child":[{"this":{}},{"tupleToUserset":{"computedUserset":{"relation":"admin"},"tupleset":{"relation":"project"}}}]}},"can_access_console":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"user"}}]}},"can_access_files":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"user"}}]}},"can_connect_sftp":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"user"}}]}},"can_edit":{"computedUserset":{"relation":"operator"}},"can_exec":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"user"}}]}},"can_manage_backups":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"operator"}}]}},"can_manage_snapshots":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"operator"}}]}},"can_update_state":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"operator"}}]}},"can_view":{"computedUserset":{"relation":"viewer"}},"operator":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"admin"}},{"tupleToUserset":{"computedUserset":{"relation":"operator"},"tupleset":{"relation":"project"}}}]}},"project":{"this":{}},"user":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"operator"}},{"tupleToUserset":{"computedUserset":{"relation":"user"},"tupleset":{"relation":"project"}}}]}},"viewer":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"user"}},{"tupleToUserset":{"computedUserset":{"relation":"viewer"},"tupleset":{"relation":"project"}}}]}}},"type":"instance"},{"metadata":{"relations":{"can_edit":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_view":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"project":{"directly_related_user_types":[{"type":"project"}]}}},"relations":{"can_edit":{"union":{"child":[{"this":{}},{"tupleToUserset":{"computedUserset":{"relation":"operator"},"tupleset":{"relation":"project"}}}]}},"can_view":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"can_edit"}},{"tupleToUserset":{"computedUserset":{"relation":"viewer"},"tupleset":{"relation":"project"}}}]}},"project":{"this":{}}},"type":"network"},{"metadata":{"relations":{"can_edit":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_view":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"project":{"directly_related_user_types":[{"type":"project"}]}}},"relations":{"can_edit":{"union":{"child":[{"this":{}},{"tupleToUserset":{"computedUserset":{"relation":"operator"},"tupleset":{"relation":"project"}}}]}},"can_view":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"can_edit"}},{"tupleToUserset":{"computedUserset":{"relation":"viewer"},"tupleset":{"relation":"project"}}}]}},"project":{"this":{}}},"type":"network_acl"},{"metadata":{"relations":{"can_edit":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_view":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"project":{"directly_related_user_types":[{"type":"project"}]}}},"relations":{"can_edit":{"union":{"child":[{"this":{}},{"tupleToUserset":{"computedUserset":{"relation":"operator"},"tupleset":{"relation":"project"}}}]}},"can_view":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"can_edit"}},{"tupleToUserset":{"computedUserset":{"relation":"viewer"},"tupleset":{"relation":"project"}}}]}},"project":{"this":{}}},"type":"network_address_set"},{"metadata":{"relations":{"can_edit":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_view":{},"server":{"directly_related_user_types":[{"type":"server"}]}}},"relations":{"can_edit":{"union":{"child":[{"this":{}},{"tupleToUserset":{"computedUserset":{"relation":"admin"},"tupleset":{"relation":"server"}}}]}},"can_view":{"tupleToUserset":{"computedUserset":{"relation":"viewer"},"tupleset":{"relation":"server"}}},"server":{"this":{}}},"type":"network_integration"},{"metadata":{"relations":{"can_edit":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_view":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"project":{"directly_related_user_types":[{"type":"project"}]}}},"relations":{"can_edit":{"union":{"child":[{"this":{}},{"tupleToUserset":{"computedUserset":{"relation":"operator"},"tupleset":{"relation":"project"}}}]}},"can_view":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"can_edit"}},{"tupleToUserset":{"computedUserset":{"relation":"viewer"},"tupleset":{"relation":"project"}}}]}},"project":{"this":{}}},"type":"network_zone"},{"metadata":{"relations":{"can_edit":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_view":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"project":{"directly_related_user_types":[{"type":"project"}]}}},"relations":{"can_edit":{"union":{"child":[{"this":{}},{"tupleToUserset":{"computedUserset":{"relation":"operator"},"tupleset":{"relation":"project"}}}]}},"can_view":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"can_edit"}},{"tupleToUserset":{"computedUserset":{"relation":"viewer"},"tupleset":{"relation":"project"}}}]}},"project":{"this":{}}},"type":"profile"},{"metadata":{"relations":{"admin":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_create_image_aliases":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_create_images":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_create_instances":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_create_network_acls":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_create_network_address_sets":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_create_network_zones":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_create_networks":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_create_profiles":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_create_storage_buckets":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_create_storage_volumes":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_edit":{},"can_view":{},"can_view_events":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_view_operations":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"operator":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"server":{"directly_related_user_types":[{"type":"server"}]},"user":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"viewer":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]}}},"relations":{"admin":{"union":{"child":[{"this":{}},{"tupleToUserset":{"computedUserset":{"relation":"admin"},"tupleset":{"relation":"server"}}}]}},"can_create_image_aliases":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"operator"}}]}},"can_create_images":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"operator"}}]}},"can_create_instances":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"operator"}}]}},"can_create_network_acls":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"operator"}}]}},"can_create_network_address_sets":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"operator"}}]}},"can_create_network_zones":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"operator"}}]}},"can_create_networks":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"operator"}}]}},"can_create_profiles":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"operator"}}]}},"can_create_storage_buckets":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"operator"}}]}},"can_create_storage_volumes":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"operator"}}]}},"can_edit":{"computedUserset":{"relation":"admin"}},"can_view":{"computedUserset":{"relation":"viewer"}},"can_view_events":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"viewer"}}]}},"can_view_operations":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"viewer"}}]}},"operator":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"admin"}},{"tupleToUserset":{"computedUserset":{"relation":"operator"},"tupleset":{"relation":"server"}}}]}},"server":{"this":{}},"user":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"operator"}},{"tupleToUserset":{"computedUserset":{"relation":"user"},"tupleset":{"relation":"server"}}}]}},"viewer":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"user"}},{"tupleToUserset":{"computedUserset":{"relation":"viewer"},"tupleset":{"relation":"server"}}}]}}},"type":"project"},{"metadata":{"relations":{"admin":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"authenticated":{"directly_related_user_types":[{"type":"user","wildcard":{}}]},"can_create_certificates":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_create_network_integrations":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_create_projects":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_create_storage_pools":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_edit":{},"can_override_cluster_target_restriction":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_view":{},"can_view_metrics":{},"can_view_privileged_events":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_view_resources":{},"can_view_sensitive":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"operator":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"user":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"viewer":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]}}},"relations":{"admin":{"this":{}},"authenticated":{"this":{}},"can_create_certificates":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"admin"}}]}},"can_create_network_integrations":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"admin"}}]}},"can_create_projects":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"admin"}}]}},"can_create_storage_pools":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"admin"}}]}},"can_edit":{"computedUserset":{"relation":"admin"}},"can_override_cluster_target_restriction":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"admin"}}]}},"can_view":{"computedUserset":{"relation":"authenticated"}},"can_view_metrics":{"computedUserset":{"relation":"authenticated"}},"can_view_privileged_events":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"admin"}}]}},"can_view_resources":{"computedUserset":{"relation":"authenticated"}},"can_view_sensitive":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"viewer"}}]}},"operator":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"admin"}}]}},"user":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"operator"}}]}},"viewer":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"user"}}]}}},"type":"server"},{"metadata":{"relations":{"can_edit":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_view":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"project":{"directly_related_user_types":[{"type":"project"}]}}},"relations":{"can_edit":{"union":{"child":[{"this":{}},{"tupleToUserset":{"computedUserset":{"relation":"operator"},"tupleset":{"relation":"project"}}}]}},"can_view":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"can_edit"}},{"tupleToUserset":{"computedUserset":{"relation":"viewer"},"tupleset":{"relation":"project"}}}]}},"project":{"this":{}}},"type":"storage_bucket"},{"metadata":{"relations":{"can_edit":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_view":{},"server":{"directly_related_user_types":[{"type":"server"}]}}},"relations":{"can_edit":{"union":{"child":[{"this":{}},{"tupleToUserset":{"computedUserset":{"relation":"admin"},"tupleset":{"relation":"server"}}}]}},"can_view":{"tupleToUserset":{"computedUserset":{"relation":"authenticated"},"tupleset":{"relation":"server"}}},"server":{"this":{}}},"type":"storage_pool"},{"metadata":{"relations":{"can_edit":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_manage_backups":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_manage_snapshots":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"can_view":{"directly_related_user_types":[{"type":"user"},{"relation":"member","type":"group"}]},"project":{"directly_related_user_types":[{"type":"project"}]}}},"relations":{"can_edit":{"union":{"child":[{"this":{}},{"tupleToUserset":{"computedUserset":{"relation":"operator"},"tupleset":{"relation":"project"}}}]}},"can_manage_backups":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"can_edit"}}]}},"can_manage_snapshots":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"can_edit"}}]}},"can_view":{"union":{"child":[{"this":{}},{"computedUserset":{"relation":"can_edit"}},{"tupleToUserset":{"computedUserset":{"relation":"viewer"},"tupleset":{"relation":"project"}}}]}},"project":{"this":{}}},"type":"storage_volume"}]}`
