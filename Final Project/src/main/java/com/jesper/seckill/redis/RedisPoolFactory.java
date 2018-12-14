package com.jesper.seckill.redis;

import java.util.LinkedHashSet;
import java.util.Set;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Bean;
import org.springframework.stereotype.Service;

import redis.clients.jedis.HostAndPort;
import redis.clients.jedis.JedisCluster;
import redis.clients.jedis.JedisPoolConfig;

@Service
public class RedisPoolFactory {

    @Autowired
    RedisConfig  redisConfig;

    /**
     * 将redis Cluster注入spring容器
     * @return
     */
    @Bean
    public JedisCluster JedisPoolFactory(){
        JedisPoolConfig config = new JedisPoolConfig();
        config.setMaxIdle(redisConfig.getPoolMaxIdle());
        config.setMaxTotal(redisConfig.getPoolMaxTotal());
        config.setMaxWaitMillis(redisConfig.getPoolMaxWait() * 1000);
        
        Set<HostAndPort> nodes = new LinkedHashSet<HostAndPort>();
    	nodes.add(new HostAndPort("127.0.0.1", 7000));
    	nodes.add(new HostAndPort("127.0.0.1", 7001));
    	nodes.add(new HostAndPort("127.0.0.1", 7002));
    	nodes.add(new HostAndPort("127.0.0.1", 7003));
    	nodes.add(new HostAndPort("127.0.0.1", 7004));
    	nodes.add(new HostAndPort("127.0.0.1", 7005));
    	nodes.add(new HostAndPort("127.0.0.1", 7006));
    	nodes.add(new HostAndPort("127.0.0.1", 7007));
    	nodes.add(new HostAndPort("127.0.0.1", 7008));
    	JedisCluster cluster = new JedisCluster(nodes, config);

        
       /* JedisPool jp = new JedisPool(config, redisConfig.getHost(), redisConfig.getPort(),
                redisConfig.getTimeout()*1000, null, 0);*/
        return cluster;
    }

}
