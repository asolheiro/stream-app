/** @type {import('next').NextConfig} */
const nextConfig = {
  experimental: {
      after: true
  },
  images: {
      remotePatterns: [
          {
            protocol: 'http',
            hostname: 'host.docker.internal',
            port: '9000',
            pathname: '/media/thumbails/**'
          },
          {
            protocol: 'http',
            hostname: 'nginx',
            port: '9000',
            pathname: '/media/thumbails/**'
          },
          {
            protocol: 'http',
            hostname: 'localhost',
            port: '9000',
            pathname: '/media/thumbails/**'
          },
          {
            protocol: 'http',
            hostname: '127.0.0.1',
            port: '9000',
            pathname: '/media/thumbails/**'
          },
          {
            protocol: 'http',
            hostname: '0.0.0.0',
            port: '9000',
            pathname: '/media/thumbails/**'
          },
          {
            protocol: 'http',
            hostname: '*',
            port: '9000',
            pathname: '/media/thumbails/**'
          }
      ]
  }
};

module.exports = nextConfig;
