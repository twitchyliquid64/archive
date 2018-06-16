import jsonrpclib

server = jsonrpclib.Server(sensor['rpc_url'])
#server.__send('RPCServer.SysStats', [1])
if sensor['component2'] == '':
    state.last_result['val'] =  server.RPCService.SysStats()[sensor['component1']]
else:
    state.last_result['val'] =  server.RPCService.SysStats()[sensor['component1']][sensor['component2']]
state.ok = True
write()
