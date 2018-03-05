import MySQLdb
import warnings

import cli


class SQLHelper(cli.CLI):
    def add_arguments(self, parser):
        parser.add_argument('--file', help='target sql script file')
        parser.add_argument('--warnings', action='store_true', help='should output warning infomation')


    def exec(self):
        if self.args.warnings:
            import sys
            warnings.filterwarnings('error', category=MySQLdb.Warning)

        with open(self.args.file) as f:
            data = f.read()
            queries = [query.strip() for query in data.split(';') if query.strip()]

        conn = MySQLdb.connect(db='lmm', user='root')
        cursor = conn.cursor()

        try:
            for query in queries:
                print(query)
                try:
                    cursor.execute(query)
                except MySQLdb.Warning as w:
                    cli.warn(w)
                print()

            return
            cli.ok('All {} sql statement(s) executed.'.format(len(queries)))

            if self.args.deploy:
                conn.commit()
                cli.info('deploy mode')
            else:
                conn.rollback()
                cli.info('dry run mode')

        except Exception as e:
            cli.error(e)
            conn.rollback()

        finally:
            cursor.close()
            conn.close()
    
    
if __name__ == '__main__':
    helper = SQLHelper()
    helper.exec()
