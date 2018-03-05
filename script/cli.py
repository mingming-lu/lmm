import sys
import argparse


white  = '\033[0m'
red    = '\033[31m'
green  = '\033[32m'
yellow = '\033[33m'


def cprint(*args, color, prefix='', **kwds):
    print(color + prefix, *args, white, **kwds)


def info(*args, **kwds):
    cprint(*args, color=white, prefix='info:', **kwds)


def ok(*args, **kwds):
    cprint(*args, color=green, prefix='ok:', **kwds)


def warn(*args, **kwds):
    cprint(*args, color=yellow, prefix='warnings:', **kwds)


def error(*args, **kwds):
    cprint(*args, color=red, prefix='error:', **kwds)


class CLI:
    def __init__(self):
        self._parse_args()

        if not self.args.deploy:
            print('##############################')
            print('######## DRYRUN START ########')
            print('##############################\n')

    def _parse_args(self):
        parser = argparse.ArgumentParser()
        parser.add_argument('--deploy', action='store_true')
        self.add_arguments(parser)
        self.args = parser.parse_args()

    def add_arguments(self):
        pass

    def exec(self, *args, **kwds):
        raise NotImplementedError

    def __del__(self):
        if not self.args.deploy:
            print('\n##############################')
            print('######### DRYRUN END #########')
            print('##############################')


if __name__ == '__main__':
    print('Please inherit CLI and implement exec')
