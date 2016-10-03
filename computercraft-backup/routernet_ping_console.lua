rednet.open("top")
while true do
	term.write("Command: ")
	local inpt = io.read()
	if inpt == "listen" then
		local ev, p1, p2, p3 = os.pullEvent("rednet_message")
		if ev == "rednet_message" then 
			print (p1, p2, p3)
		end
	end
	if inpt == "ping" then
		term.write("RouterID of subnet to ping: ")
		local ID = io.read()
		rednet.broadcast("SERVMSG~"..ID.."~CONSOLE~PING")
	end
	if inpt == "consoleserv" then
		while true do
			local ev, p1, p2, p3 = os.pullEvent("rednet_message")
			if ev == "rednet_message" then 
				print (p1, p2, p3)
			end
		end
	end
end
