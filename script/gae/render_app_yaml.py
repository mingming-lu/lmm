import argparse
import sys


def parse_args():
    parser = argparse.ArgumentParser(description='default', argument_default='')
    parser.add_argument('--input')
    parser.add_argument('--output')

    known, unknown = parser.parse_known_args()
    extra = _parse_attribute(unknown)

    return vars(known), vars(extra)


def _parse_attribute(args):
    parser = argparse.ArgumentParser(description='attribute')
    for arg in args:
        if arg.startswith(('-', '--')):
            parser.add_argument(arg)

    return parser.parse_args(args)


if __name__ == '__main__':
    knowns, unkowns = parse_args()

    with open(knowns['input']) as f:
        template = f.read()

    yaml_file = template.format(**unkowns)

    with open(knowns['output'], 'w') as f:
        f.write(yaml_file)
