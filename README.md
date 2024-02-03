ip tunnel add tun-6in4 mode sit remote <client ipv4 addr> local <local ipv4> 
ip link set tun-6in4  up
ip addr add 2001:db8:::1/64 dev tun-6in4

ip route add 

ip tunnel add tun-6in4 mode sit remote 1.6.4.1 local 185.185.58.180
# list of commands
ip tunnel add tun-6in4 mode sit remote 185.60.45.135 local 185.185.58.180
ip addr add 2a00::1/127 dev tun-6in4

# tunnel up
ip link set tun-6in4 up

# tunnel down
ip link set tun-6in4 down

# add route via interface
ip -6 route add 2a00:1::/64 dev tun-6in4

# check ipv6 route via tunnel interface
ip -6 route show


# template
ip tunnel add {{tun-name}} mode sit remote {{remote-ipv4}} local {{local-ipv4}}
ip link set {{tun-name}} up
ip addr add {{local-ipv6}} dev {{tun-name}}
ip -6 route add {{ipv6-pd}} dev {{tun-name}}



# API description

1. register the client
1.1 Receive  email and   password
1.2 generate new api key

2. create 6in4 tunnel
2.1 define nearest server for the client
2.2 get client's ipv4 address
2.3 get ipv6 interface addresses /127
2.4 get ipv6 pd prefix /64
2.5 create template for interface
2.6 provisioning the interface and route via interface

3. update the tunnel
3.1 receive update with proper API key, update the ipv4 remote address.
    unconfigure->configure new tunnel

4. For registered users
    list of configured tunnels
        actions: Add, Delete
    API key:
        actions: Update
