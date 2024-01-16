import argparse
import os
import time

from sklearn.metrics import mean_squared_error

import joblib
import pandas as pd
import vineyard

def test_model(with_vineyard):
    os.system('sync; echo 3 > /proc/sys/vm/drop_caches')
    st = time.time()
    if with_vineyard:
        client = vineyard.connect(config="/vineyard/client/config")
        x_test_data = client.get(name="/data/x_test.pkl", fetch=True)
        y_test_data = client.get(name="/data/y_test.pkl", fetch=True)
    else:
        x_test_data = pd.read_pickle("/data/x_test.pkl")
        y_test_data = pd.read_pickle("/data/y_test.pkl")
        #delete the x_test.pkl and y_test.pkl
        os.remove("/data/x_test.pkl")
        os.remove("/data/y_test.pkl")
    ed = time.time()
    print('##################################')
    print('read x_test and y_test execution time: ', ed - st)

    model = joblib.load("/data/model.pkl")
    y_pred = model.predict(x_test_data)

    err = mean_squared_error(y_test_data, y_pred)

    with open('/data/output.txt', 'a') as f:
        f.write(str(err))


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('--with_vineyard', type=bool, default=False, help='Whether to use vineyard')
    args = parser.parse_args()
    st = time.time()
    print('Testing model...')
    test_model(args.with_vineyard)
    ed = time.time()
    print('##################################')
    print('Testing model data time: ', ed - st)
