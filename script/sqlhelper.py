import MySQLdb
import warnings

import cli


class SQLHelper(cli.CLI):
    '''
    usage: sqlhelper.py [-h] [--deploy] [--file FILE] [--warnings]
    '''

    def add_arguments(self, parser):
        parser.add_argument('--file', help='target sql script file')
        parser.add_argument('--database', help='target database')
        parser.add_argument('--warnings', action='store_true', help='should output warning infomation')

    def exec(self):
        if self.args.warnings:
            warnings.filterwarnings('error', category=MySQLdb.Warning)
        else:
            warnings.filterwarnings('ignore', category=MySQLdb.Warning)

        try:
            with open(self.args.file) as f:
                data = f.read()
                queries = [query.strip() for query in data.split(';') if query.strip()]
        except Exception as e:
            cli.error(e)
            return

        try:
            db = MySQLdb.connect(user="root", host='lmm-mysql')
            cursor = db.cursor()
            cursor.execute('CREATE DATABASE IF NOT EXISTS {};'.format(self.args.database))
            db.select_db(self.args.database)
            for query in queries:
                print(query)
                try:
                    cursor.execute(query)
                except MySQLdb.Warning as w:
                    cli.warn(w)
                print()

            cli.ok('All {} sql statement(s) executed.'.format(len(queries)))

            if self.args.deploy:
                db.commit()
            else:
                db.rollback()

        except Exception as e:
            cli.error(e)
            db.rollback()

        finally:
            cursor.close()
            db.close()
    
    
if __name__ == '__main__':
    helper = SQLHelper()
    helper.exec()
