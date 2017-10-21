-- Container object is the representation of a Docker
-- container in the Minecraft world

-- NewMonitoringStatusContainer returns a Container object,
-- representation of a container in
-- the Minecraft world
function NewMonitoringStatusContainer()
  c = {
    displayed = false,
    x = 0,
    z = 0,
    name="",
    id="",
    percent=0,
    init=MonitoringStatusContainer.init,
    setInfos=MonitoringStatusContainer.setInfos,
    destroy=MonitoringStatusContainer.destroy,
    display=MonitoringStatusContainer.display,
    addGround=MonitoringStatusContainer.addGround
  }
  return c
end

MonitoringStatusContainer = {
    displayed = false, 
    x = 0,
    z = 0, 
    name="",
    id="",
    percent=0,
}

-- MonitoringStatusContainer:init sets Container's position
function MonitoringStatusContainer:init(x,z)
  self.x = x
  self.z = z
  self.displayed = false
end

-- MonitoringStatusContainer:setInfos sets Container's id, name, imageRepo,
-- image tag and running state
function MonitoringStatusContainer:setInfos(id,name,percent)
  self.id = id
  self.name = name
  self.percent = percent
end

-- MonitoringStatusContainer:display displays all Container's blocks
-- Blocks will be blue if the container is running,
-- orange otherwise.
function MonitoringStatusContainer:display()

  local metaPrimaryColor = E_META_WOOL_LIGHTBLUE
  local metaSecondaryColor = E_META_WOOL_ORANGE

  self.displayed = true

  local counter = self.percent
  for i=1, 10
  do
    if counter > 0
    then
      counter = counter - 10
      setBlock(UpdateQueue,self.x+1,GROUND_TABLE_LEVEL+i,self.z + 1,E_BLOCK_WOOL,metaSecondaryColor) 
    else
      setBlock(UpdateQueue,self.x+1,GROUND_TABLE_LEVEL+i,self.z + 1,E_BLOCK_WOOL,metaPrimaryColor)
    end
  end
end

-- MonitoringStatusContainer:addGround creates ground blocks
-- necessary to display the container
function MonitoringStatusContainer:addGround()
  local y = GROUND_TABLE_LEVEL
  local max_x = GROUND_TABLE_MAX_X

  if GROUND_TABLE_MIN_X > self.x - 2
  then
    max_x = GROUND_TABLE_MIN_X
    GROUND_TABLE_MIN_X = self.x - 2
    min_x = GROUND_TABLE_MIN_X
  end

  local min_x = GROUND_TABLE_MIN_X
  for x= min_x, max_x
  do
    for z=GROUND_TABLE_MIN_Z,GROUND_TABLE_MAX_Z
    do
      setBlock(UpdateQueue,x,y,z,E_BLOCK_GRASS,0)
      for sky=y+1,y+6
      do
        setBlock(UpdateQueue,x,sky,z,E_BLOCK_AIR,0)
      end
    end
  end
end
