#!/bin/bash
# test_network_address_set.sh
# A standalone test suite for Incus network address sets.
# This script exercises creation (CLI, API, from STDIN, with project scoping),
# listing, show, edit, patch, adding/removing addresses, renaming, custom keys, and deletion.
#
# Requirements:
#   - The "incus" CLI must be installed and in PATH.
#   - Optionally, INCUS_ADDR can be set (defaults to http://localhost:8443).
#
# I made this because I was unable to run test in the standard way so I made a workaround 
#

set -euo pipefail

# --- Helpers for colored output ---
function info() {
  echo -e "\033[1;34m[INFO]\033[0m $1"
}
function success() {
  echo -e "\033[1;32m[SUCCESS]\033[0m $1"
}
function error_msg() {
  echo -e "\033[1;31m[ERROR]\033[0m $1"
}

# --- Helper functions ---
function get_container_ip() {
  local container="$1"
  local ip=""
  for i in {1..10}; do
    ip=$(incus list "$container" --format csv | cut -d',' -f3 | head -n1 | cut -d' ' -f1)
    if [ -n "$ip" ]; then break; fi
    sleep 1
  done
  echo "$ip"
}

function get_container_ip6() {
  local container="$1"
  local ip6=""
  for i in {1..10}; do
    ip6=$(incus list "$container" --format csv | cut -d',' -f4 | tr ' ' '\n' | head -n1)
    if [ -n "$ip6" ]; then break; fi
    sleep 1
  done
  echo "$ip6"
}


get_container_ip_ovn() { incus list testct --format json | jq -r '.[0].state.network.eth0.addresses[] | select(.family=="inet").address'; }
get_container_ip_ovn6() { incus list testct --format json | jq -r '.[0].state.network.eth0.addresses[] | select(.family=="inet6").address' | head -n 1; }

# --- Check that incus is installed ---
if ! command -v incus &> /dev/null; then
  error_msg "incus CLI could not be found. Please install it first."
  exit 1
fi

# --- Set default INCUS_ADDR if not provided ---
export INCUS_ADDR="${INCUS_ADDR:-http://localhost:8443}"



# --- Test functions ---

# Test 1: Creation using CLI (rejecting non‐hostname compatible names, then creating and deleting)
function test_creation_cli() {
  info "Test 1: Creation (CLI)"
  if incus network address-set create 2432; then
    error_msg "Non-hostname compatible name was accepted, expected rejection."
    exit 1
  else
    success "Non-hostname compatible name correctly rejected."
  fi

  incus network address-set create testAS
  success "Address set 'testAS' created."
  incus network address-set delete testAS
  success "Address set 'testAS' deleted."
}

# Test 2: Creation using project scoping
function test_creation_project() {
  info "Test 2: Creation (Project)"
  incus project create testproj -c features.networks=true
  incus network address-set create testAS --project testproj
  if incus network address-set ls --project testproj | grep -q "testAS"; then
    success "Address set 'testAS' exists in project 'testproj'."
  else
    error_msg "Address set 'testAS' not found in project 'testproj'."
    exit 1
  fi
  # Clean up
  incus network address-set delete testAS --project testproj
  incus project delete testproj
  success "Project 'testproj' and its address set cleaned up."
}

# Test 3: Creation from STDIN
function test_creation_stdin() {
  info "Test 3: Creation from STDIN"
  # Note: Ensure addresses is a YAML list (not a plain string)
  cat <<EOF | incus network address-set create testAS
description: Test Address set from STDIN
addresses:
  - 192.168.0.1
  - 192.168.0.254
external_ids:
  user.mykey: foo
EOF
  if incus network address-set show testAS | grep -q "description: Test Address set from STDIN"; then
    success "Address set created from STDIN with correct description."
  else
    error_msg "Failed: Address set created from STDIN did not have the expected description."
    exit 1
  fi
  incus network address-set delete testAS
}

# Test 4: Listing address sets
function test_listing() {
  info "Test 4: Listing"
  incus network address-set create testAS --description "Listing test"
  if incus network address-set ls | grep -q "testAS"; then
    success "Address set 'testAS' appears in listing."
  else
    error_msg "Address set 'testAS' does not appear in listing."
    exit 1
  fi
  incus network address-set delete testAS
}

# Test 5: Show command
function test_show() {
  info "Test 5: Show"
  incus network address-set create testAS --description "Show test"
  if incus network address-set show testAS | grep -q "description: Show test"; then
    success "Show command returns correct description."
  else
    error_msg "Show command did not return expected description."
    exit 1
  fi
  incus network address-set delete testAS
}

# Test 6: Edit command (using STDIN)
function test_edit() {
  info "Test 6: Edit"
  incus network address-set create testAS --description "Initial description"
  cat <<EOF | incus network address-set edit testAS
description: Updated address set
addresses:
  - 10.0.0.1
  - 10.0.0.2
external_ids:
  user.mykey: bar
EOF
  if incus network address-set show testAS | grep -q "Updated address set"; then
    success "Edit command updated the description correctly."
  else
    error_msg "Edit command failed to update the description."
    exit 1
  fi
  incus network address-set delete testAS
}

# Test 7: Patch command (update external IDs only)
function test_patch() {
  info "Test 7: Patch"
  incus network address-set create testAS --description "Patch test"
  incus query -X PATCH -d "{\"external_ids\": {\"user.myotherkey\": \"bah\"}}" /1.0/network-address-sets/testAS
  if incus network address-set show testAS | grep -q "user.myotherkey: bah"; then
    success "Patch command updated external IDs correctly."
  else
    error_msg "Patch command did not update external IDs as expected."
    exit 1
  fi
  incus network address-set delete testAS
}

# Test 8: Add and Remove addresses
function test_add_remove_addresses() {
  info "Test 8: Add/Remove Addresses"
  incus network address-set create testAS --description "Address add/remove test"
  # Add an address using the CLI subcommand "add-addr"
  incus network address-set add-addr testAS 192.168.1.100
  if incus network address-set show testAS | grep -q "192.168.1.100"; then
    success "Address 192.168.1.100 added."
  else
    error_msg "Failed to add address 192.168.1.100."
    exit 1
  fi
  incus network address-set remove-addr testAS 192.168.1.100
  if ! incus network address-set show testAS | grep -q "192.168.1.100"; then
    success "Address 192.168.1.100 removed."
  else
    error_msg "Failed to remove address 192.168.1.100."
    exit 1
  fi
  incus network address-set delete testAS
}

# Test 9: Rename command
function test_rename() {
  info "Test 9: Rename"
  incus network address-set create testAS --description "Rename test"
  incus network address-set rename testAS testAS-renamed
  if incus network address-set ls | grep -q "testAS-renamed"; then
    success "Rename succeeded: testAS-renamed found in listing."
  else
    error_msg "Rename failed: testAS-renamed not found."
    exit 1
  fi
  incus network address-set delete testAS-renamed
}

# Test 10: Custom keys (set/unset)
function test_custom_keys() {
  info "Test 10: Custom Keys"
  incus network address-set create testAS --description "Custom keys test"
  incus network address-set set testAS user.somekey foo
  if incus network address-set show testAS | grep -q "foo"; then
    success "Custom key 'user.somekey' set to foo."
  else
    error_msg "Failed to set custom key 'user.somekey'."
    exit 1
  fi
  incus network address-set unset testAS user.somekey
  if ! incus network address-set show testAS | grep -q "foo"; then
    success "Custom key 'user.somekey' successfully unset."
  else
    error_msg "Custom key 'user.somekey' was not unset."
    exit 1
  fi
  incus network address-set delete testAS
}

# Test 11: Deletion
function test_delete() {
  info "Test 11: Deletion"
  incus network address-set create testAS --description "Delete test"
  incus network address-set delete testAS
  if incus network address-set ls | grep -q "testAS"; then
    error_msg "Address set 'testAS' still exists after deletion."
    exit 1
  else
    success "Address set 'testAS' successfully deleted."
  fi
}

# Test 12: Block ping using address-sets

function test_nft_block_ping_with_address_set() {
    info "Test 12: ACL block ICMP for container"
    local ip=$(get_container_ip testct)
    info "Container IPv4: $ip"
    if ping -c2 "$ip" > /dev/null; then
        success "Ping to container succeeded."
    else
        error_msg "Ping to container failed."
        incus delete testct --force
        exit 1
    fi
    incus network address-set create testAS
    incus network address-set add-addr testAS $ip
    incus network acl create blockping
    incus network acl rule add blockping ingress action=drop protocol=icmp4 destination="\$testAS"
    incus network set incusbr0 security.acls="blockping"
    sleep 2
    if ping -c2 "$ip" > /dev/null; then
        error_msg "Ping succeeded despite ACL block."
        incus network set incusbr0 security.acls=""
        incus network acl delete blockping
        incus delete testct --force
        exit 1
    else
        success "Ping correctly blocked by ACL."
    fi
    incus network address-set remove-addr testAS $ip
    if ping -c2 "$ip" > /dev/null; then
        success "Ping to container succeeded."
    else
        error_msg "Ping to container failed."
        incus delete testct --force
        exit 1
    fi
    incus network set incusbr0 security.acls=""
    incus network acl delete blockping
    incus network address-set delete testAS
}

# Test 13 Block pingv6 using address-sets

function test_nft_block_pingv6_with_address_set() {
    info "Test 13: ACL block ICMPv6 for container"
    local ip=$(get_container_ip6 testct)
    info "Container IPv6: $ip"
    if ping -c2 "$ip" > /dev/null; then
        success "Ping to container succeeded."
    else
        error_msg "Ping to container failed."
        #incus delete testct --force
        exit 1
    fi
    incus network address-set create testAS
    incus network address-set add-addr testAS $ip
    incus network acl create blockping
    incus network acl rule add blockping ingress action=drop protocol=icmp6 destination="\$testAS"
    incus network set incusbr0 security.acls="blockping"
    sleep 2
    if ping -c2 "$ip" > /dev/null; then
        error_msg "Ping succeeded despite ACL block."
        #incus network set incusbr0 security.acls=""
        #incus network acl delete blockping
        #incus delete testct --force
        exit 1
    else
        success "Ping correctly blocked by ACL."
    fi
    incus network address-set remove-addr testAS $ip
    if ping -c2 "$ip" > /dev/null; then
        success "Ping to container succeeded."
    else
        error_msg "Ping to container failed."
        incus delete testct --force
        exit 1
    fi
    incus network set incusbr0 security.acls=""
    incus network acl delete blockping
    incus network address-set rm testAS
}

# Test 14

function test_nft_acl_mixed_subject() {
    info "Test 14: ACL with mixed subject (literal IP and address set)"
    # Create ACL that blocks TCP port 22 if destination is either a literal IP or an address set.
    incus network address-set create testAS
    local ip=$(get_container_ip testct)
    incus network address-set add-addr testAS "$ip"
    incus launch images:debian/12 testct2
    sleep 3
    local ip2=$(get_container_ip testct2)
    incus network acl create mixedACL
    incus network acl rule add mixedACL ingress action=drop protocol=icmp4 destination="$ip2,\$testAS"
    incus network set incusbr0 security.acls="mixedACL"
    sleep 2
    if ping -c2 "$ip" > /dev/null; then
        error_msg "Ping succeeded despite ACL block; expected failure."
        incus network set incusbr0 security.acls=""
        incus network acl delete mixedACL
        incus network address-set delete testAS
        incus delete testct2 --force
        incus delete testct --force
        exit 1
    else
        success "Ping correctly blocked by mixed ACL."
    fi
    if ping -c2 "$ip2" > /dev/null; then
        error_msg "Ping succeeded despite ACL block; expected failure."
        incus network set incusbr0 security.acls=""
        incus network acl delete mixedACL
        incus network address-set delete testAS
        incus delete testct2 --force
        incus delete testct --force
        exit 1
    else
        success "Ping correctly blocked by mixed ACL."
    fi
    incus network set incusbr0 security.acls=""
    incus network acl delete mixedACL
    incus network address-set rm testAS
    incus delete testct2 --force
}

# Test 15

function test_nft_update_with_cidr() {
    info "Test 15: Update address set with container network CIDR and verify ACL block"
    local ip=$(get_container_ip testct)
    local subnet=$(echo "$ip" | awk -F. '{print $1"."$2"."$3".0/24"}')
    info "Derived subnet: $subnet"
    incus network address-set create testAS
    incus network address-set add-addr testAS "$subnet"
    incus network acl create cidrACL
    incus network acl rule add cidrACL ingress action=drop protocol=icmp4 destination="\$testAS"
    incus network set incusbr0 security.acls="cidrACL"
    sleep 2
    if ping -c2 "$ip" > /dev/null; then
    error_msg "Ping succeeded despite CIDR ACL block; expected failure."
    incus network set incusbr0 security.acls=""
    incus network acl delete cidrACL
    incus network address-set delete testAS
    incus delete testct --force
    exit 1
    else
    success "Ping correctly blocked by CIDR ACL."
    fi
    incus network set incusbr0 security.acls=""
    incus network acl delete cidrACL
    incus network address-set rm testAS
}

# Test 16 

function test_nft_block_tcp(){
    incus exec testct -- sh -c "apt install netcat-openbsd -y"
    local ip=$(get_container_ip testct)
    incus network address-set create testAS
    incus network address-set add-addr testAS "$ip"
    incus network acl create blocktcp7896
    incus network acl rule add blocktcp7896 ingress action=drop protocol=tcp destination_port="7896" destination="\$testAS"
    incus network set incusbr0 security.acls="blocktcp7896"
    echo "Please run 'incus exec testct -- sh -c \"nohup nc -l -p 7896 >/dev/null 2>&1\"' Press ENTER when done"
    read _
    if nc -z -w 5 "$ip" 7896; then
      error_msg "TCP connection to port 7896 on testct succeeded despite ACL block; expected failure."
      incus network set incusbr0 security.acls=""
      incus network acl delete blocktcp7896
      incus network address-set rm testAS
      incus delete testct --force
      exit 1
    else
      success "TCP connection to port 7896 on testct IPv4 correctly blocked by ACL."
    fi
    incus network set incusbr0 security.acls=""
    incus network acl delete blocktcp7896
    incus network address-set rm testAS
}

# Test 17

function test_nft_block_tcp_mixed_set(){
    incus exec testct -- sh -c "apt install netcat-openbsd -y"
    local ip=$(get_container_ip testct)
    local ip6=$(get_container_ip6 testct)
    incus network address-set create testAS
    incus network address-set add-addr testAS "$ip"
    incus network address-set add-addr testAS "$ip6"
    incus network acl create blocktcp7896
    incus network acl rule add blocktcp7896 ingress action=drop protocol=tcp destination_port="7896" destination="\$testAS"
    incus network set incusbr0 security.acls="blocktcp7896"
    echo "Please run 'incus exec testct -- sh -c \"nohup nc -l -p 7896 >/dev/null 2>&1\"' Press ENTER when done"
    read _
    if nc -z -w 5 "$ip" 7896; then
      error_msg "TCP connection to port 7896 on testct succeeded despite ACL block; expected failure."
      incus network set incusbr0 security.acls=""
      incus network acl delete blocktcp7896
      incus network address-set rm testAS
      incus delete testct --force
      exit 1
    else
      success "TCP connection to port 7896 on testct IPv6 correctly blocked by ACL."
    fi
    echo "Please run 'incus exec testct -- sh -c \"nohup nc -6 -l -p 7896 >/dev/null 2>&1\"' Press ENTER when done"
    read _
    if nc -6 -z -w 5 "$ip6" 7896; then
      error_msg "TCP connection to port 7896 on testct succeeded despite ACL block; expected failure."
      incus network set incusbr0 security.acls=""
      incus network acl delete blocktcp7896
      incus network address-set rm testAS
      incus delete testct --force
      exit 1
    else
      success "TCP connection to port 7896 on testct IPv6 correctly blocked by ACL."
    fi
    incus network set incusbr0 security.acls=""
    incus network acl delete blocktcp7896
    incus network address-set rm testAS
}

# Test 18

function test_nft_dynamic_ping_block() {
    info "Test: Dynamically updating an address set for ACL blocking"

    # Step 1: Get container IPs
    local ip_testct=$(get_container_ip testct)
    
    incus launch images:debian/12 testct2
    sleep 3
    local ip_testct2=$(get_container_ip testct2)
    info "Container testct IPv4: $ip_testct"
    info "Container testct2 IPv4: $ip_testct2"

    # Step 2: Create address set and add both IPs
    incus network address-set create testAS
    incus network address-set add-addr testAS "$ip_testct"
    incus network address-set add-addr testAS "$ip_testct2"

    # Step 3: Create ACL to block ICMP to testAS
    incus network acl create blockping
    incus network acl rule add blockping ingress action=drop protocol=icmp4 destination="\$testAS"
    incus network set "incusbr0" security.acls="blockping"

    sleep 2

    # Step 4: Verify that both containers are blocked
    if ping -c2 "$ip_testct" > /dev/null; then
        error_msg "Ping to testct succeeded despite ACL block."
        incus network acl delete blockping
        incus network address-set delete testAS
        exit 1
    else
        success "Ping to testct correctly blocked."
    fi

    if ping -c2 "$ip_testct2" > /dev/null; then
        error_msg "Ping to testct2 succeeded despite ACL block."
        incus network acl delete blockping
        incus network address-set delete testAS
        exit 1
    else
        success "Ping to testct2 correctly blocked."
    fi

    # Step 5: Remove testct2's IP from the address set
    incus network address-set remove-addr testAS "$ip_testct2"
    sleep 1

    # Step 6: Verify that testct is still blocked but testct2 is now reachable
    if ping -c2 "$ip_testct" > /dev/null; then
        error_msg "Ping to testct succeeded despite ACL block."
        incus network acl delete blockping
        incus network address-set delete testAS
        exit 1
    else
        success "Ping to testct correctly remains blocked."
    fi

    if ping -c2 "$ip_testct2" > /dev/null; then
        success "Ping to testct2 succeeded after address set update."
    else
        error_msg "Ping to testct2 failed despite being removed from the blocked set."
        incus network acl delete blockping
        incus network address-set delete testAS
        exit 1
    fi
    incus network set incusbr0 security.acls=""
    incus network acl delete blockping
    incus network address-set delete testAS
    incus rm --force testct2
}

### OVN TESTS

### TEST 19: Block ICMPv4 using address-sets
function test_ovn_block_ping_with_address_set() {
    info "Test 12: ACL block ICMPv4 for container"
    local ip=$(get_container_ip_ovn)
    info "Container IPv4: $ip"

    if ping -c2 "$ip" > /dev/null; then
        success "Ping to container succeeded."
    else
        error_msg "Ping to container failed."
    fi

    incus network address-set create testAS
    incus network address-set add-addr testAS "$ip"
    incus network acl create blockping
    incus network acl rule add blockping ingress action=drop protocol=icmp4 destination="\$testAS"
    incus network set "$OVN_NETWORK" security.acls="blockping"
    sleep 2

    if ping -c2 "$ip" > /dev/null; then
        error_msg "Ping succeeded despite ACL block."
    else
        success "Ping correctly blocked by ACL."
    fi

    incus network set "$OVN_NETWORK" security.acls=""
    incus network acl delete blockping
    incus network address-set delete testAS
}

### TEST 20: Block ICMPv6 using address-sets
function test_ovn_block_pingv6_with_address_set() {
    info "Test 13: ACL block ICMPv6 for container"
    local ip=$(get_container_ip_ovn6)
    info "Container IPv6: $ip"

    if ping -6 -c2 "$ip" > /dev/null; then
        success "Ping to container succeeded."
    else
        error_msg "Ping to container failed."
    fi

    incus network address-set create testAS
    incus network address-set add-addr testAS "$ip"
    incus network acl create blockping
    incus network acl rule add blockping ingress action=drop protocol=icmp6 destination="\$testAS"
    incus network set "$OVN_NETWORK" security.acls="blockping"

    sleep 2
    if ping -6 -c2 "$ip" > /dev/null; then
        error_msg "Ping succeeded despite ACL block."
    else
        success "Ping correctly blocked by ACL."
    fi

    incus network set "$OVN_NETWORK" security.acls=""
    incus network acl delete blockping
    incus network address-set delete testAS
}

### TEST 21: Mixed ACL Subject
function test_ovn_acl_mixed_subject() {
    info "Test 14: ACL with mixed subject (literal IP and address set)"
    local ip=$(get_container_ip_ovn)

    incus network address-set create testAS
    incus network address-set add-addr testAS "$ip"
    incus launch images:debian/12 testct2
    sleep 3
    local ip2=$(get_container_ip_ovn testct2)

    incus network acl create mixedACL
    incus network acl rule add mixedACL ingress action=drop protocol=icmp4 destination="$ip2,\$testAS"
    incus network set "$OVN_NETWORK" security.acls="mixedACL"

    sleep 2
    if ping -c2 "$ip" > /dev/null || ping -c2 "$ip2" > /dev/null; then
        error_msg "Ping succeeded despite ACL block."
    else
        success "Ping correctly blocked by ACL."
    fi

    incus network set "$OVN_NETWORK" security.acls=""
    incus network acl delete mixedACL
    incus network address-set delete testAS
    incus delete testct2 --force
}

### TEST 22: CIDR Address Set
function test_ovn_update_with_cidr() {
    info "Test 15: Update address set with CIDR"
    local ip=$(get_container_ip_ovn)
    local subnet=$(echo "$ip" | awk -F. '{print $1"."$2"."$3".0/24"}')

    incus network address-set create testAS
    incus network address-set add-addr testAS "$subnet"
    incus network acl create cidrACL
    incus network acl rule add cidrACL ingress action=drop protocol=icmp4 destination="\$testAS"
    incus network set "$OVN_NETWORK" security.acls="cidrACL"

    sleep 2
    if ping -c2 "$ip" > /dev/null; then
        error_msg "Ping succeeded despite CIDR ACL block."
    else
        success "Ping correctly blocked by ACL."
    fi

    incus network set "$OVN_NETWORK" security.acls=""
    incus network acl delete cidrACL
    incus network address-set delete testAS
}

### TEST 23: IPv4/TCP Block
function test_ovn_block_tcp() {
    local ip=$(get_container_ip_ovn)
    incus exec testct -- apt update && apt install -y netcat-openbsd

    incus network address-set create testAS
    echo "Container ip: $ip"
    incus network address-set add-addr testAS "$ip"
    incus network acl create blocktcp7896
    incus network acl rule add blocktcp7896 ingress action=drop protocol=tcp destination_port="7896" destination="\$testAS"
    incus network set "$OVN_NETWORK" security.acls="blocktcp7896"

    info "Run: incus exec testct -- nohup nc -l -p 7896 &"
    read

    if nc -z -w 5 "$ip" 7896; then
        error_msg "TCP connection succeeded despite ACL block."
    else
        success "TCP connection correctly blocked by ACL."
    fi

    incus network set "$OVN_NETWORK" security.acls=""
    incus network acl delete blocktcp7896
    incus network address-set delete testAS
}

### TEST 24: TCP Block with Mixed Address Sets
function test_ovn_block_tcp_mixed_set() {
    local ip=$(get_container_ip_ovn)
    local ip6=$(get_container_ip_ovn6)
    incus exec testct -- apt update && apt install -y netcat-openbsd

    incus network address-set create testAS
    incus network address-set add-addr testAS "$ip"
    incus network address-set add-addr testAS "$ip6"
    incus network acl create blocktcp7896
    incus network acl rule add blocktcp7896 ingress action=drop protocol=tcp destination_port="7896" destination="\$testAS"
    incus network set "$OVN_NETWORK" security.acls="blocktcp7896"

    info "Run: incus exec testct -- nohup nc -l -p 7896 &"
    read

    if nc -z -w 5 "$ip" 7896 ; then
        error_msg "TCP connection succeeded despite ACL block."
    else
        success "TCP connection correctly blocked by ACL."
    fi
    info "Run: incus exec testct -- nohup nc -6 -l -p 7896 &"
    read
    if nc -6 -z -w 5 "$ip6" 7896; then
        error_msg "TCP connection succeeded despite ACL block."
    else
        success "TCP connection correctly blocked by ACL."
    fi

    incus network set "$OVN_NETWORK" security.acls=""
    incus network acl delete blocktcp7896
    incus network address-set delete testAS
}

### TEST 25

function test_ovn_dynamic_ping_block() {
    info "Test: Dynamically updating an address set for ACL blocking"

    # Step 1: Get container IPs
    local ip_testct=$(get_container_ip_ovn testct)
    local ip_testct2=$(get_container_ip_ovn testct2)
    incus init images:debian/12 testct2
    incus config device override testct2 eth0 network=ovntest
    incus start testct2
    sleep 3
    info "Container testct IPv4: $ip_testct"
    info "Container testct2 IPv4: $ip_testct2"

    # Step 2: Create address set and add both IPs
    incus network address-set create testAS
    incus network address-set add-addr testAS "$ip_testct"
    incus network address-set add-addr testAS "$ip_testct2"

    # Step 3: Create ACL to block ICMP to testAS
    incus network acl create blockping
    incus network acl rule add blockping ingress action=drop protocol=icmp4 destination="\$testAS"
    incus network set "$OVN_NETWORK" security.acls="blockping"

    sleep 2

    # Step 4: Verify that both containers are blocked
    if ping -c2 "$ip_testct" > /dev/null; then
        error_msg "Ping to testct succeeded despite ACL block."
        incus network acl delete blockping
        incus network address-set delete testAS
        exit 1
    else
        success "Ping to testct correctly blocked."
    fi

    if ping -c2 "$ip_testct2" > /dev/null; then
        error_msg "Ping to testct2 succeeded despite ACL block."
        incus network acl delete blockping
        incus network address-set delete testAS
        exit 1
    else
        success "Ping to testct2 correctly blocked."
    fi

    # Step 5: Remove testct2's IP from the address set
    incus network address-set remove-addr testAS "$ip_testct2"
    sleep 2

    # Step 6: Verify that testct is still blocked but testct2 is now reachable
    if ping -c2 "$ip_testct" > /dev/null; then
        error_msg "Ping to testct succeeded despite ACL block."
        incus network acl delete blockping
        incus network address-set delete testAS
        exit 1
    else
        success "Ping to testct correctly remains blocked."
    fi

    if ping -c2 "$ip_testct2" > /dev/null; then
        success "Ping to testct2 succeeded after address set update."
    else
        error_msg "Ping to testct2 failed despite being removed from the blocked set."
        incus network acl delete blockping
        incus network address-set delete testAS
        exit 1
    fi

    incus network acl delete blockping
    incus network address-set delete testAS
    incus rm --force testct2
}


## --- Run all tests ---
info "Starting network address set tests..."
#test_creation_cli
#success "TEST 1 CLI CREATION OK"
#test_creation_project
#success "TEST 2 CREATION PROJECT OK"
#test_creation_stdin
#success "TEST 3 CREATION STDIN OK"
#test_listing
#success "TEST 4 LISTING OK"
#test_show
#success "TEST 5 SHOW OK"
#test_edit
#success "TEST 6 EDIT OK"
#test_patch
#success "TEST 7 PATCH OK"
#test_add_remove_addresses
#success "TEST 8 REMOVE ADDR OK"
#test_rename
#success "TEST 9 RENAME OK"
#test_custom_keys
#success "TEST 10 CUSTOM KEYS OK"
#test_delete
#success "TEST 11 DELETION OK"

#info Tests needing a container begin now
incus launch images:debian/12 testct
sleep 3

info Testing nftables behaviour

#test_nft_block_ping_with_address_set
#success "TEST 12 BLOCK ICMPv4 OK"
#test_nft_block_pingv6_with_address_set # NOK
#success "TEST 13 BLOCK ICMPv6 OK"
#test_nft_acl_mixed_subject
#success "TEST 14 MIXED ACL SUBJECTS OK"
#test_nft_update_with_cidr
#success "TEST 15 CIDR ADDRESS SET OK"
test_nft_block_tcp
success "TEST 16 IPv4/TCP BLOCK OK"
#test_nft_block_tcp_mixed_set
#success "TEST 17 TCP BLOCK MIXED ADDRESS SET OK"
#test_nft_dynamic_ping_block
#success "TEST 18 NFT DYNAMIC PING BLOCK OK"
#incus rm --force testct

info "OVN tests begins"

# Ensure the container network is OVN-based
PARENT_NETWORK="incusbr0"
OVN_NETWORK="ovntest"

info "Ensuring OVN network is set up..."
incus network set "$PARENT_NETWORK" ipv4.dhcp.ranges="10.158.174.100-10.158.174.110" ipv4.ovn.ranges="10.158.174.111-10.158.174.120"
incus network create "$OVN_NETWORK" --type=ovn network="$PARENT_NETWORK" || true
# Get ovn ipv4 external address and setup routing:
ovnnet=$(incus network show ovntest | grep 'ipv4.address' | head -n 1 | cut -d ':' -f2)
ovnip=$(incus network show ovntest | grep 'network.ipv4.address' | head -n 1 | cut -d ':' -f2)
# Add route to ovn network
ip r a $(echo $ovnnet | awk -F. '{print $1"."$2"."$3".0/24"}') via $ovnip
# Allow icmp to go through ovn network
lsin=$(ovn-nbctl ls-list | awk '/-ls-int/ {gsub(/[()]/, "", $2); print $2}')
ovn-nbctl acl-add $lsin to-lport 200 "(icmp4 || icmp6)" allow
# Start a test container
info "Launching test container..."
incus stop --force testct
sleep 3
incus config device override testct eth0 network=ovntest
incus start testct
sleep 3

#test_ovn_block_ping_with_address_set
#success "TEST 19 OVN PING BLOCK OK"
#test_ovn_block_pingv6_with_address_set  # NOK
#success "TEST 20 OVN PING6 OK"
#test_ovn_acl_mixed_subject
#success "TEST 21 OVN MIXED ACL SUBJECT OK"
#test_ovn_update_with_cidr
#success "TEST 22 OVN CIDR ADDRESS SET OK"
#test_ovn_block_tcp

success "TEST 23 OVN IPv4/TCP BLOCK OK"
test_ovn_block_tcp_mixed_set
#success "TEST 24 OVN TCP BLOCK MIXED ADDRESS SET OK"
#test_ovn_dynamic_ping_block
#success "TEST 25 OVN DYNAMIC PING BLOCK OK"

incus delete --force testct
incus network rm ovntest