import pandas as pd
from prophet import Prophet
import json

def lambda_handler(event, context):
    print(event)
    body = json.loads(event['body'])

    df = pd.read_json(json.dumps(body['input']))

    m = Prophet(changepoint_prior_scale=body['changepoint_prior_scale'], interval_width=body['interval_width'])
    m.fit(df)

    future = m.make_future_dataframe(periods=1, freq="h")

    forecast = m.predict(future)
    return {
        'statusCode': 200,
        'body' : forecast[['ds', 'yhat', 'yhat_lower', 'yhat_upper']].to_json()
    }

# result = lambda_handler({'body': '{"changepointPriorScale":0.25,"intervalWidth":0.95,"input":{"ds":["2024-10-08T17:40:00","2024-10-08T18:10:00","2024-10-08T18:40:00","2024-10-08T19:10:00","2024-10-08T19:40:00","2024-10-08T20:10:00","2024-10-08T20:40:00","2024-10-08T21:10:00","2024-10-08T21:40:00","2024-10-08T22:10:00","2024-10-08T22:40:00","2024-10-08T23:10:00","2024-10-08T23:40:00","2024-10-09T00:10:00","2024-10-09T00:40:00","2024-10-09T01:10:00","2024-10-09T01:40:00","2024-10-09T02:10:00","2024-10-09T02:40:00","2024-10-09T03:10:00","2024-10-09T03:40:00","2024-10-09T04:10:00","2024-10-09T04:40:00","2024-10-09T05:10:00","2024-10-09T05:40:00","2024-10-09T06:10:00","2024-10-09T06:40:00","2024-10-09T07:10:00","2024-10-09T07:40:00","2024-10-09T08:10:00","2024-10-09T08:40:00","2024-10-09T09:10:00","2024-10-09T09:40:00","2024-10-09T10:10:00","2024-10-09T10:40:00","2024-10-09T11:10:00","2024-10-09T11:40:00","2024-10-09T12:10:00","2024-10-09T12:40:00","2024-10-09T13:10:00","2024-10-09T13:40:00","2024-10-09T14:10:00","2024-10-09T14:40:00","2024-10-09T15:10:00","2024-10-09T15:40:00","2024-10-09T16:10:00","2024-10-09T16:40:00","2024-10-09T17:10:00","2024-10-09T17:40:00","2024-10-09T18:10:00","2024-10-09T18:40:00","2024-10-09T19:10:00","2024-10-09T19:40:00","2024-10-09T20:10:00","2024-10-09T20:40:00","2024-10-09T21:10:00","2024-10-09T21:40:00","2024-10-09T22:10:00","2024-10-09T22:40:00","2024-10-09T23:10:00","2024-10-09T23:40:00","2024-10-10T00:10:00","2024-10-10T00:40:00","2024-10-10T01:10:00","2024-10-10T01:40:00","2024-10-10T02:10:00","2024-10-10T02:40:00","2024-10-10T03:10:00","2024-10-10T03:40:00","2024-10-10T04:10:00","2024-10-10T04:40:00","2024-10-10T05:10:00","2024-10-10T05:40:00","2024-10-10T06:10:00","2024-10-10T06:40:00","2024-10-10T07:10:00","2024-10-10T07:40:00","2024-10-10T08:10:00","2024-10-10T08:40:00","2024-10-10T09:10:00","2024-10-10T09:40:00","2024-10-10T10:10:00","2024-10-10T10:40:00","2024-10-10T11:10:00","2024-10-10T11:40:00","2024-10-10T12:10:00","2024-10-10T12:40:00","2024-10-10T13:10:00","2024-10-10T13:40:00","2024-10-10T14:10:00","2024-10-10T14:40:00","2024-10-10T15:10:00","2024-10-10T15:40:00","2024-10-10T16:10:00","2024-10-10T16:40:00","2024-10-10T17:10:00","2024-10-10T17:40:00","2024-10-10T18:10:00","2024-10-10T18:40:00","2024-10-10T19:10:00"],"y":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,2769,634,585,559,822,512,503,553,542,539,514,589,502,544,501,554,2911,4336,4240,4289,4289,4321,4328,4332,1964,502,521,507,475,503,540,542,526,530,497,568,520,567,688,1387,1544,495,203,0,0,0,0,0,0,0,0,0,0,0,99,1658,3310,4933]}}'}, {})

# print(result)