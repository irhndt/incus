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
  incus network address-set create testAS --description "Patch test"
  incus query -X PATCH -d "{\"external_ids\": {\"user.myotherkey\": \"bah\"}}" /1.0/network-address-sets/testAS
  incus network address-set show testAS | grep -q "user.myotherkey: bah"
  incus network address-set delete testAS
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
}
