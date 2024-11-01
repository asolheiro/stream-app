import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  experimental: {
    after: true
  },
  images : {
    remotePatterns: [
      {
        hostname: 'localhost'
      },
      {
        hostname: 'host.docker.internal'
      }
    ]
  }
};

export default nextConfig;
