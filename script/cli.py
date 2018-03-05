import sys
import argparse


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
