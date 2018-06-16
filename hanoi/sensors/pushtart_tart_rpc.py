import jsonrpclib

server = jsonrpclib.Server(sensor['rpc_url'])
state.last_result['val'] = server.Tarts.GetTart(APIKey=sensor['api_key'],PushURL=sensor['push_url'])['Tart']
if 'component1' in sensor:
    state.last_result['val'] = state.last_result['val'][sensor['component1']]
if 'component2' in sensor:
    state.last_result['val'] = state.last_result['val'][sensor['component2']]
if 'component3' in sensor:
    state.last_result['val'] = state.last_result['val'][sensor['component3']]
state.ok = True
write()
