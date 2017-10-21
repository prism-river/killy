StatusContainers = {}
EmptyContainerSpace = {}

-- updateContainer accepts 3 different states: running, stopped, created
-- sometimes "start" events arrive before "create" ones
-- in this case, we just ignore the update
function updateStatusContainer(id,name,percent)
  LOG("Update container with ID: " .. id .. " state: " .. tostring(percent))

  -- first pass, to see if the container is
  -- already displayed (maybe with another state)
  for i=1, table.getn(StatusContainers)
  do
    -- if container found with same ID, we update it
    if StatusContainers[i] ~= EmptyContainerSpace and StatusContainers[i].id == id
    then
      StatusContainers[i]:setInfos(id,name,percent)
      StatusContainers[i]:display()
      LOG("found. updated. now return")
      return
    end
  end

  -- if container isn't already displayed, we see if there's an empty space
  -- in the world to display the container
  local x = STATUS_CONTAINER_START_X
  local index = -1

  for i=1, table.getn(StatusContainers)
  do
    -- use first empty location
    if StatusContainers[i] == EmptyContainerSpace
    then
      LOG("Found empty location: StatusContainers[" .. tostring(i) .. "]")
      index = i
      break
    end
    x = x + ACTIVE_CONTAINER_OFFSET_X
  end

  LOG("create a new active statuss container")
  local container = NewMonitoringStatusContainer()
  container:init(x,STATUS_CONTAINER_START_Z)
  container:setInfos(id,name,percent)
  container:addGround()
  container:display()

  if index == -1
  then
    table.insert(StatusContainers, container)
  else
    StatusContainers[index] = container
  end
end
