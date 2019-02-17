import argparse
import re

class Server:
    @classmethod
    def new_from_string(cls, name: str):
        matched = re.match(r'(\w+)(:(\d+))?', name)
        hostname, _, port, *_ = matched.groups()
        return Server(hostname, port)

    def __init__(self, hostname: str, port: int):
        self._hostname = hostname
        self._port     = port

    def __str__(self):
        if self._port is not None:
            return f'{self._hostname}:{self._port}'
        return self._hostname


class Upstream:
    def __init__(self, name: str):
        self._name    = name
        self._servers = []

    def name(self) -> str:
        return self._name

    def add_server(self, server: Server):
        self._servers.append(server)

    def remove_server(self, name):
        if name not in self._servers:
            print(f'[warn] no such server ({name}) found in upstream: {self._name}')
        self._servers.remove(name)

    def servers(self):
        return self._servers

    def __iter__(self):
        return iter(self._servers)


class Config:
    def __init__(self, file):
        self._upstreams = {}
        self._parse_config_file(file)

    def _parse_config_file(self, file):
        upstream_begin  = re.compile(r'upstream\s+(.+)\s+{')
        upstream_end    = re.compile('}')
        upstream_server = re.compile(r'\s+server\s+(\w+)(:(\d+))?')

        current_upstream = None

        for line in file:
            if not line.strip():
                continue
            res = upstream_begin.match(line)
            if res:
                name, *_ = res.groups()
                upstream = Upstream(name)
                self._upstreams[upstream.name()] = upstream
                current_upstream = upstream
                continue

            res = upstream_server.match(line)
            if res:
                hostname, _, port, *_ = res.groups()
                current_upstream.add_server(Server(hostname, port))
                continue

            if upstream_end.match(line):
                current_upstream = None
                continue

    def write_to_file(self, dst):
        with open(dst, mode='w') as f:
            for upstream in self:
                f.write(f'upstream {upstream.name()} {{\n')
                for server in upstream:
                    f.write(f'    server {server} fail_timeout=2s;\n')
                f.write('}\n\n')

    def upstream(self, name):
        return self._upstreams[name]

    def add_upstream(self, upstreams):
        self._upstreams[upstream.name()] = upstream

    def remove_upstream(self, name):
        if name not in self._upstreams:
            print(f'[warn] no such upstream found: {name}')
        del self._upstreams[name]

    def __contains__(self, upstream: str) -> bool:
        return upstream in self._upstreams

    def __iter__(self):
        return iter(self._upstreams.values())

    def __str__(self):
        pass

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument('--src', '-s', help='src NGINX upstream config file')
    parser.add_argument('--dst', '-d', help='dst NGINX upstream config file')
    parser.add_argument('--add_upstream', action='append')
    parser.add_argument('--rm_upstream', action='append')
    parser.add_argument('--reset_upstream', action='append')

    args = parser.parse_args()

    with open(args.src) as f:
        conf = Config(f)

    if args.add_upstream:
        for value in args.add_upstream:
            upstream_name, *servers = value.split(' ')
            if upstream_name in conf:
                for server in servers:
                    conf.upstream(upstream_name).add_server(Server.new_from_string(server))
            else:
                upstream = Upstream(upstream_name)
                for server in servers:
                    upstream.add_server(Server.new_from_string(server))
                conf.add_upstream(upstream)

    if args.rm_upstream:
        for value in args.rm_upstream:
            conf.remove_upstream(value)

    if args.reset_upstream:
        for value in args.reset_upstream:
            upstream_name, *servers = value.split(' ')
            upstream = Upstream(upstream_name)
            for server in servers:
                upstream.add_server(Server.new_from_string(server))
            conf.add_upstream(upstream)

    conf.write_to_file(args.dst)
