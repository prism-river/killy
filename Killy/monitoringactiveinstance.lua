ActiveInstanceContainers = {}
EmptyContainerSpace = {}

-- destroyContainer looks for the first container having the given id,
-- removes it from the Minecraft world and from the 'ActiveInstanceContainers' array
function destroyContainer(id)
  LOG("destroyContainer: " .. id)
  -- loop over the containers and remove the first having the given id
  for i=1, table.getn(ActiveInstanceContainers)
  do
    if ActiveInstanceContainers[i] ~= EmptyContainerSpace and ActiveInstanceContainers[i].id == id
    then
      -- remove the container from the world
      ActiveInstanceContainers[i]:destroy()
      -- if the container being removed is the last element of the array
      -- we reduce the size of the "Container" array, but if it is not,
      -- we store a reference to the "EmptyContainerSpace" object at the
      -- same index to indicate this is a free space now.
      -- We use a reference to this object because it is not possible to
      -- have 'nil' values in the middle of a lua array.
      if i == table.getn(ActiveInstanceContainers)
      then
        table.remove(ActiveInstanceContainers, i)
        -- we have removed the last element of the array. If the array
        -- has tailing empty container spaces, we remove them as well.
        while ActiveInstanceContainers[table.getn(ActiveInstanceContainers)] == EmptyContainerSpace
        do
          table.remove(ActiveInstanceContainers, table.getn(ActiveInstanceContainers))
        end
      else
        ActiveInstanceContainers[i] = EmptyContainerSpace
      end
      -- we removed the container, we can exit the loop
      break
    end
  end
end

-- updateContainer accepts 3 different states: running, stopped, created
-- sometimes "start" events arrive before "create" ones
-- in this case, we just ignore the update
function updateActiveInstanceContainer(id,name,state)
  LOG("Update container with ID: " .. id .. " state: " .. tostring(state))

  -- first pass, to see if the container is
  -- already displayed (maybe with another state)
  for i=1, table.getn(ActiveInstanceContainers)
  do
    -- if container found with same ID, we update it
    if ActiveInstanceContainers[i] ~= EmptyContainerSpace and ActiveInstanceContainers[i].id == id
    then
      ActiveInstanceContainers[i]:setInfos(id,name,state)
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
  container:setInfos(id,name,true)
  container:addGround()
  container:display(true)

  if index == -1
  then
    table.insert(ActiveInstanceContainers, container)
  else
    ActiveInstanceContainers[index] = container
  end
end
