import argparse


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('--APP_ORIGIN')
    parser.add_argument('--MANAGER_ORIGIN')
    parser.add_argument('--ASSET_BUCKET_NAME')
    parser.add_argument('--DATASTORE_PROJECT_ID')
    parser.add_argument('--GIN_MODE')
    parser.add_argument('--GOOGLE_APPLICATION_CREDENTIALS')
    parser.add_argument('--LMM_API_TOKEN_KEY')
    parser.add_argument('--YAML_FILE', default='app.yaml')

    args = parser.parse_args()

    with open('script/template.yaml') as f:
        template = f.read()

    yaml_file = template.format(
        APP_ORIGIN=args.APP_ORIGIN,
        MANAGER_ORIGIN=args.MANAGER_ORIGIN,
        ASSET_BUCKET_NAME=args.ASSET_BUCKET_NAME,
        DATASTORE_PROJECT_ID=args.DATASTORE_PROJECT_ID,
        GIN_MODE=args.GIN_MODE,
        GOOGLE_APPLICATION_CREDENTIALS=args.GOOGLE_APPLICATION_CREDENTIALS,
        LMM_API_TOKEN_KEY=args.LMM_API_TOKEN_KEY,
    )

    with open(args.YAML_FILE, 'w') as f:
        f.write(yaml_file)
