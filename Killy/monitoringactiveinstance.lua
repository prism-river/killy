ActiveInstanceContainers = {}
EmptyContainerSpace = {}

-- updateContainer accepts 3 different states: running, stopped, created
-- sometimes "start" events arrive before "create" ones
-- in this case, we just ignore the update
function updateActiveInstanceContainer(id,name,realname,state)
  LOG("Update container with ID: " .. id .. " state: " .. tostring(state))

  -- first pass, to see if the container is
  -- already displayed (maybe with another state)
  for i=1, table.getn(ActiveInstanceContainers)
  do
    -- if container found with same ID, we update it
    if ActiveInstanceContainers[i] ~= EmptyContainerSpace and ActiveInstanceContainers[i].id == id
    then
      ActiveInstanceContainers[i]:setInfos(id,name,realname,state)
      ActiveInstanceContainers[i]:display(state)
      LOG("found. updated. now return")
      return
    end
  end

  -- if container isn't already displayed, we see if there's an empty space
  -- in the world to display the container
  local x = ACTIVE_CONTAINER_START_X
  local index = -1

  for i=1, table.getn(ActiveInstanceContainers)
  do
    -- use first empty location
    if ActiveInstanceContainers[i] == EmptyContainerSpace
    then
      LOG("Found empty location: ActiveInstanceContainers[" .. tostring(i) .. "]")
      index = i
      break
    end
    x = x + ACTIVE_CONTAINER_OFFSET_X
  end

  LOG("create a new active server container")
  local container = NewActiveInstanceContainer()
  container:init(x,ACTIVE_CONTAINER_START_Z)
  container:setInfos(id,name,realname,true)
  container:addGround()
  container:display(true)

  if index == -1
  then
    table.insert(ActiveInstanceContainers, container)
  else
    ActiveInstanceContainers[index] = container
  end
end

-- updateStats update CPU and memory usage displayed
-- on container sign (container identified by id)
function updateActiveInstanceStats(id, cap, a)
  for i=1, table.getn(ActiveInstanceContainers)
  do
    if ActiveInstanceContainers[i] ~= EmptyContainerSpace and ActiveInstanceContainers[i].id == id
    then
      ActiveInstanceContainers[i]:updateCapacitySign(cap)
      ActiveInstanceContainers[i]:updateAvailabilitySign(a)
      break
    end
  end
end
