TableRecords = {}
EmptyContainerSpace = {}

-- updateContainer accepts 3 different states: running, stopped, created
-- sometimes "start" events arrive before "create" ones
-- in this case, we just ignore the update
function updateTableRecordContainer(id,name,content)
  LOG("Update tablet record container with ID: " .. id .. " content: " .. content)
  -- first pass, to see if the container is
  -- already displayed (maybe with another state)
  for i=1, table.getn(TableRecords)
  do
    -- if container found with same ID, we update it
    if TableRecords[i] ~= EmptyContainerSpace and TableRecords[i].id == id
    then
      TableRecords[i]:setInfos(id,name,content)
      TableRecords[i]:display()
      LOG("found. updated. now return")
      return
    end
  end

  -- if container isn't already displayed, we see if there's an empty space
  -- in the world to display the container
  local x = TABLE_AREA_START_X
  local index = -1

  for i=1, table.getn(TableRecords)
  do
    -- use first empty location
    if TableRecords[i] == EmptyContainerSpace
    then
      LOG("Found empty location: Containers[" .. tostring(i) .. "]")
      index = i
      break
    end
    x = x + TABLE_SIGNAL_OFFSET
  end

  LOG("create a new tablet record container")
  local container = NewTableRecordContainer()
  container:init(x,CONTAINER_START_Z)
  container:setInfos(id,name,content)
  container:addGround()
  container:display()

  if index == -1
  then
    table.insert(TableRecords, container)
  else
    TableRecords[index] = container
  end
end
