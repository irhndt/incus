# shellcheck disable=2148
test_address_set() {
  ensure_import_testimage
  ensure_has_localhost_remote "${INCUS_ADDR}"
  ! incus network address-set create 2432 || false
  incus network address-set create testAS
  incus network address-set delete testAS
  # Test 2: Address set creation & deletion
  ! incus network address-set create 2432 || false
  incus network address-set create testAS
  incus network address-set delete testAS
  incus project create testproj -c features.networks=true
  incus network address-set create testAS --project testproj
  incus network address-set ls --project testproj | grep -q "testAS"
  incus network address-set delete testAS --project testproj
  incus project delete testproj
  cat <<EOF | incus network address-set create testAS
description: Test Address set from STDIN
addresses:
  - 192.168.0.1
  - 192.168.0.254
external_ids:
  user.mykey: foo
EOF
  incus network address-set show testAS | grep -q "description: Test Address set from STDIN"
  incus network address-set delete testAS
  incus network address-set create testAS --description "Listing test"
  incus network address-set ls | grep -q "testAS"
  incus network address-set delete testAS
  incus network address-set create testAS --description "Show test"
  incus network address-set delete testAS
  incus network address-set create testAS --description "Initial description"
  cat <<EOF | incus network address-set edit testAS
description: Updated address set
addresses:
  - 10.0.0.1
  - 10.0.0.2
external_ids:
  user.mykey: bar
EOF
  incus network address-set show testAS | grep -q "Updated address set"
  incus network address-set delete testAS
  # This was OK on my machine
  # incus network address-set create testAS --description "Patch test"
  # incus query -X PATCH -d "{\"external_ids\": {\"user.myotherkey\": \"bah\"}}" /1.0/network-address-sets/testAS
  # incus network address-set show testAS | grep -q "user.myotherkey: bah"
  # incus network address-set delete testAS
  incus network address-set create testAS --description "Address add/remove test"
  incus network address-set add-addr testAS 192.168.1.100
  incus network address-set show testAS | grep -q "192.168.1.100"
  incus network address-set del-addr testAS 192.168.1.100
  ! incus network address-set show testAS | grep -q "192.168.1.100" || false
  incus network address-set delete testAS
  incus network address-set create testAS --description "Rename test"
  incus network address-set rename testAS testAS-renamed
  incus network address-set ls | grep -q "testAS-renamed"
  incus network address-set delete testAS-renamed
  incus network address-set create testAS --description "Custom keys test"
  incus network address-set set testAS user.somekey foo
  incus network address-set show testAS | grep -q "foo"
  incus network address-set delete testAS
  ! incus network address-set ls | grep -q "testAS" || false
  # Testing if address sets are working correctly inside acls
  brName="inct$$"

  # Standard bridge.
  incus network create "${brName}" \
        ipv6.dhcp.stateful=true \
        ipv4.address=192.0.2.1/24 \
        ipv6.address=2001:db8::1/64

  incus launch images:debian/12 testct
  ip=$(incus list testct --format csv | cut -d',' -f3 | head -n1 | cut -d' ' -f1)
  incus network address-set create testAS
  incus network address-set add-addr testAS "$ip"
  incus network acl create allowping
  incus network acl rule add blockping ingress action=allow protocol=icmp4 destination="\$testAS"
  incus network set "${brName}" security.acls="allowping"
  sleep 1
  ping -c2 "$ip" > /dev/null
  incus network address-set del-addr testAS "$ip"
  incus network set "${brName}" security.acls=""
  incus network acl delete allowping
  incus network address-set delete testAS
  incus network address-set create testAS
  incus network address-set add-addr testAS "$ip"
  incus launch images:debian/12 testct2
  ip2=$(incus list testct2 --format csv | cut -d',' -f3 | head -n1 | cut -d' ' -f1)
  incus network acl create mixedACL
  incus network acl rule add mixedACL ingress action=allow protocol=icmp4 destination="$ip2,\$testAS"
  incus network set "${brName}" security.acls="mixedACL"
  sleep 1
  ping -c2 "$ip" > /dev/null
  ping -c2 "$ip2" > /dev/null
  incus network set "${brName}" security.acls=""
  incus network acl delete mixedACL
  incus network address-set rm testAS
  incus delete testct2 --force
  subnet=$(echo "$ip" | awk -F. '{print $1"."$2"."$3".0/24"}')
  incus network address-set create testAS
  incus network address-set add-addr testAS "$subnet"
  incus network acl create cidrACL
  incus network acl rule add cidrACL ingress action=allow protocol=icmp4 destination="\$testAS"
  incus network set "${brName}" security.acls="cidrACL"
  sleep 1
  ping -c2 "$ip" > /dev/null
  incus network set "${brName}" security.acls=""
  incus network acl delete cidrACL
  incus network address-set rm testAS
  ip6=$(incus list testct --format csv | cut -d',' -f4 | tr ' ' '\n' | head -n1)
  socat - TCP:"$ip":5355 # SHOULD WORK BY DEFAULT
  socat - TCP6:"$ip6":5355 5355 # SHOULD WORK BY DEFAULT
  incus network address-set create testAS
  incus network address-set add-addr testAS "$ip"
  incus network acl create allowtcp5355
  incus network acl rule add allowtcp5355 ingress action=allow protocol=tcp destination_port="5355" destination="\$testAS"
  incus network set "${brName}" security.acls="allowtcp5355"
  socat - TCP:"$ip":5355
  incus network address-set add-addr testAS "$ip6"
  socat - TCP6:"$ip6":5355 5355
  incus network address-set del-addr testAS "$ip6"
  ! socat - TCP6:"$ip6":5355 || false
  incus network set "${brName}" security.acls=""
  incus network acl delete allowtcp5355
  incus network address-set rm testAS
}
