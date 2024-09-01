#define _DEFAULT_SOURCE
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/tcp.h>
#include <linux/udp.h>
#include <linux/bpf_helpers.h>


#define SEC(NAME) __attribute__((section(NAME), used))

struct bpf_map_def SEC("maps") blocked_ips = {
    .type = BPF_MAP_TYPE_HASH,
    .key_size = sizeof(__u32),
    .value_size = sizeof(__u8),
    .max_entries = 1024,
};

struct cachestat_range {
    uint64_t off;
    uint64_t len;
};

struct cachestat {
    uint64_t nr_cache;
    uint64_t nr_dirty;
    uint64_t nr_writeback;
    uint64_t nr_evicted;
    uint64_t nr_recently_evicted;
};

SEC("xdp")
int xdp_firewall_prog(struct __sk_buff *skb) {
    void *data = (void *)(long)skb->data;
    void *data_end = (void *)(long)skb->data_end;
    struct ethhdr *eth = data;

    // Ensure packet is large enough for Ethernet header
    if ((void *)(eth + 1) > data_end)
        return XDP_PASS;

    // Check if it's an IPv4 packet
    if (eth->h_proto == htons(ETH_P_IP)) {
        struct iphdr *ip = (struct iphdr *)(eth + 1);

        // Ensure packet is large enough for IP header
        if ((void *)(ip + 1) > data_end)
            return XDP_PASS;

        // Check if the source IP is blocked
        __u32 src_ip = ip->saddr;
        __u8 *blocked = bpf_map_lookup_elem(&blocked_ips, &src_ip);
        if (blocked) {
            // Drop the packet if the IP is blocked
            return XDP_DROP;
        }
    }

    return XDP_PASS;
}

char _license[] SEC("license") = "GPL";
