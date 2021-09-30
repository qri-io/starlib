load("dataframe.star", "dataframe")


csv_text = """latitude,longitude,depth,mag,mag_type,nst,gap,dmin,rms,net,id,updated,place,type,horizontal_error,depth_error,mag_error,mag_nst,status,location_source,mag_source
31.6893333,-114.5405,9.98,5.49,mw,13,106,0.3208,0.35,ci,ci38385946,2020-03-07T21:25:58.100Z,71km SE of Estacion Coahuila,earthquake,0.78,31.61,,6,reviewed,ci,cis"""


def f():
  df = dataframe.parse_csv(csv_text)
  print(df)
  print('')

  df = df.drop(columns=['latitude'])
  print(df)
  print('')

  df = df.drop(columns=['longitude'])
  print(df)
  print('')

  df = df.drop(columns=['location_source'])
  print(df)
  print('')

  df = dataframe.DataFrame([['alpha%d' % (n+1), chr(n+65)] for n in range(26)],
                           columns=['id', 'char'])
  print(df)
  print('')


f()
