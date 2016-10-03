
--top: autoarchic gate on wooden pipe (for pumping out disks back to ender chest)
--right: ender chest (used for exchanging disks between routers)
--diagonal top-right: floppy drive (connected by wired modem)
--
--Please start with no disk in the disk drive.
--Top 10 or so slots in ender chest should be for spare disks. Only assign routerIDs after that.
--router ids are assigned by creating the file router_config and writing a single number in there.
--This number corresponds to a slot in the ender chest, which is used to send disks with data in them
--to an endpoint.

local exchangeChest = peripheral.wrap("right")
local diskDrive
local routerID

local devices = peripheral.getNames()
for subscript, side in pairs(devices) do --iterates through connected devices to locate the modem(s) and the disk drive. ONE disk drive assumed.
  if peripheral.getType(side) == "drive" then
    diskDrive = side
  end
  if peripheral.getType(side) == "modem" then
    rednet.open(side)
  end
end

configHandle = fs.open("router_config", "r")
if configHandle == nil then
	error("No configuration file present! (router_config)")
	shell.exit()
end

idStr = configHandle.readAll()
configHandle.close()
routerID = tonumber(idStr)

print("Router configured to use ID: ".. tostring(routerID))
if routerID < 10 then
	error("RouterID invalid! (ID must be between 10-21)")
	shell.exit()
end


function pushBlankDrive () --finds a spare drive and pushes it into the disk drive.
  while true do
	  for i=0,10 do
		local dataIn = exchangeChest.getStackInSlot(i)
		if dataIn ~= nil then
			if dataIn["id"] == 4257 then
				exchangeChest.push("up", i, 1)
				return
			end
		end
	  end
  print("No free drives! send() operation stalling!")
  os.sleep(3)
  end
end


function readDiskFile(fname) --given a file name, reads the file name on the disk drive and returns the contents as a string.
	local path = fs.combine(disk.getMountPath(diskDrive), fname)
	h = fs.open(path, "r")
	if h == nil then
		return ""
	end
	local tmp = h.readAll()
	h.close()
	return tmp
end

function writeDiskFile(fname, data) --given a file name and data, (over)writes a file on disk with that name with the contents of 'data' as a string.
	local path = fs.combine(disk.getMountPath(diskDrive), fname)
	h = fs.open(path, "w")
	h.write(data)
	h.close()
end

function parseMessages()
	local count = tonumber(readDiskFile("msgcount")) --count the number of messages already on the disk.
	for i=1,count do
		if readDiskFile("intent"..tostring(i)) == "unicastmessage" then
			rednet.send(tonumber(readDiskFile("destid"..tostring(i))), "PACKET~"..readDiskFile("sendingendpoint"..tostring(i)).."~"..readDiskFile("sendingcomputerid"..tostring(i)).."~"..readDiskFile("data"..tostring(i)))
		end
		if readDiskFile("intent"..tostring(i)) == "servermessage" then
			rednet.broadcast("SERVERMESSAGE~".. readDiskFile("destserver"..tostring(i)) .."~"..readDiskFile("sendingendpoint"..tostring(i)).."~"..readDiskFile("sendingcomputerid"..tostring(i)).."~"..readDiskFile("data"..tostring(i)))
		end
	end
end


function sendMessage(message, destID, destEndpoint, sendingComputerID, sendingEndpoint) --sends a UNICAST message to a machine on a specific subnet/endpoint/router/whateveryouwanttocallit.
	if exchangeChest.getStackInSlot(tonumber(destEndpoint)) == nil then --get a disk from the inbound slot if there is one and append to that, else get a blank disk.
		pushBlankDrive () --get a drive
	else
		exchangeChest.pushIntoSlot("up", tonumber(destEndpoint), 1, 0)
	end
	local count = tonumber(readDiskFile("msgcount")) --count the number of messages already on the disk.
	if count == nil then --setup and increment the count variable.
		count = 1
	else
		count = count + 1
	end
	print("Now writing message to disk. Existing messages: " .. tostring(count))
	writeDiskFile("msgcount", tostring(count))
	writeDiskFile("intent" .. tostring(count), "unicastmessage")
	writeDiskFile("destid" .. tostring(count), tostring(destID))
	writeDiskFile("destendpoint" .. tostring(count), tostring(destEndpoint))
	writeDiskFile("sendingcomputerid" .. tostring(count), tostring(sendingComputerID))
	writeDiskFile("sendingendpoint" .. tostring(count), tostring(sendingEndpoint))
	writeDiskFile("data" .. tostring(count), tostring(message))
	while exchangeChest.getStackInSlot(tonumber(destEndpoint)) ~= nil do --wait till the slot is empty (should be) to insert the disk again.
		print("sendMessage() Stalled! Waiting for reciever slot to be empty!")
		os.sleep(2)
	end
	exchangeChest.pullIntoSlot("up", 0, 1, tonumber(destEndpoint))
end


function sendServerMessage(message, servername, destEndpoint, sendingComputerID, sendingEndpoint) --sends a UNICAST message to a named server on a specific subnet/endpoint/router/whateveryouwanttocallit.
	if exchangeChest.getStackInSlot(tonumber(destEndpoint)) == nil then --get a disk from the inbound slot if there is one and append to that, else get a blank disk.
		pushBlankDrive () --get a drive
	else
		exchangeChest.pushIntoSlot("up", tonumber(destEndpoint), 1, 0)
	end
	local count = tonumber(readDiskFile("msgcount")) --count the number of messages already on the disk.
	if count == nil then --setup and increment the count variable.
		count = 1
	else
		count = count + 1
	end
	print("Now writing message to disk. Existing messages: " .. tostring(count))
	writeDiskFile("msgcount", tostring(count))
	writeDiskFile("intent" .. tostring(count), "servermessage")
	writeDiskFile("destserver" .. tostring(count), tostring(servername))
	writeDiskFile("destendpoint" .. tostring(count), tostring(destEndpoint))
	writeDiskFile("sendingcomputerid" .. tostring(count), tostring(sendingComputerID))
	writeDiskFile("sendingendpoint" .. tostring(count), tostring(sendingEndpoint))
	writeDiskFile("data" .. tostring(count), tostring(message))
	while exchangeChest.getStackInSlot(tonumber(destEndpoint)) ~= nil do --wait till the slot is empty (should be) to insert the disk again.
		print("sendMessage() Stalled! Waiting for reciever slot to be empty!")
		os.sleep(2)
	end
	exchangeChest.pullIntoSlot("up", 0, 1, tonumber(destEndpoint))
end


function wipeAndReturnDisk() --instructs pipe above computer to remove the disk from the drive, after wiping its contents.
	local files = fs.list( disk.getMountPath(diskDrive) ) --# all the files and folders in the directory '/' --> main
	for i = 1, #files do
		if files[i] ~= "rom" then --# You cannot delete rom
		   fs.delete( fs.combine(disk.getMountPath(diskDrive), files[i]) )
		end
	end

	redstone.setOutput("top", true)
	os.sleep(1)
	redstone.setOutput("top", false)
end

function split(str, pat) --utility function used to split elements of the string up so a message can be processed.
  local t = { }
  local fpat = "(.-)"..pat
  local last_end = 1
  local s, e, cap = str:find(fpat, 1)
  while s do
    if s ~= 1 or cap ~= "" then
      table.insert(t,cap)
    end
    last_end = e+1
    s, e, cap = str:find(fpat, last_end)
  end
  if last_end <= #str then
    cap = str:sub(last_end)
    table.insert(t, cap)
  end
  return t
end

os.startTimer(5) --checks for inbound messages every 5 seconds
while true do --mainloop
	local ev, p1, p2, p3 = os.pullEvent()
	if ev == "timer" then --if time to check for message
		local dataIn = exchangeChest.getStackInSlot(routerID) --see if disk (representing inbound data) is in OUR slot
		if dataIn ~= nil then
			if dataIn["id"] == 4257 then --floppy disk id 
				exchangeChest.push("up", routerID, 1) --take disk out and put it in the disk drive.
				parseMessages()
				wipeAndReturnDisk()
			end
		end
		os.startTimer(5)
	end
	if ev == "rednet_message" then --event is a new rednet message - lets see if we should relay it.
		local elements = split(p2, "~")
		if #elements ~= 4 then
			error("Protocol Error!: " .. p2)
		else
			if elements[1] == "UNICAST" then
				sendMessage(elements[4], elements[3], elements[2], tostring(p1), tostring(routerID))
			end
			if elements[1] == "SERVMSG" then
				sendServerMessage(elements[4], elements[3], elements[2], tostring(p1), tostring(routerID))
			end
		end
	end
end
