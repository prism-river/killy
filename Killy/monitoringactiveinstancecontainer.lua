-- ActiveInstanceContainer object is the representation of a Docker
-- container in the Minecraft world

-- constant variables
-- CONTAINER_CREATED = 0
-- CONTAINER_RUNNING = 1
-- CONTAINER_STOPPED = 2

-- NewActiveInstanceContainer returns a ActiveInstanceContainer object,
-- representation of a Docker container in
-- the Minecraft world
function NewActiveInstanceContainer()
  c = {
    displayed = false,
    x = 0,
    z = 0,
    name="",
    id="",
    realname="",
    running=false,
    init=ActiveInstanceContainer.init,
    setInfos=ActiveInstanceContainer.setInfos,
    destroy=ActiveInstanceContainer.destroy,
    display=ActiveInstanceContainer.display,
    updateCapacitySign=ActiveInstanceContainer.updateCapacitySign,
    updateAvailabilitySign=ActiveInstanceContainer.updateAvailabilitySign,
    addGround=ActiveInstanceContainer.addGround
  }
  return c
end

ActiveInstanceContainer = {displayed = false, x = 0, z = 0, name="",id="",realname="",imageRepo="",imageTag="",running=false}

-- ActiveInstanceContainer:init sets ActiveInstanceContainer's position
function ActiveInstanceContainer:init(x,z)
  self.x = x
  self.z = z
  self.displayed = false
end

-- ActiveInstanceContainer:setInfos sets ActiveInstanceContainer's id, name, imageRepo,
-- image tag and running state
function ActiveInstanceContainer:setInfos(id,name,realname,running)
  self.id = id
  self.name = name
  self.realname = realname
  self.running = running
end

-- ActiveInstanceContainer:display displays all ActiveInstanceContainer's blocks
-- Blocks will be blue if the container is running,
-- orange otherwise.
function ActiveInstanceContainer:display(running)

  local metaPrimaryColor = E_META_WOOL_LIGHTBLUE
  local metaSecondaryColor = E_META_WOOL_BLUE

  if running == false
  then
    metaPrimaryColor = E_META_WOOL_ORANGE
    metaSecondaryColor = E_META_WOOL_RED
  end

  self.displayed = true

  for px=self.x, self.x+3
  do
    for pz=self.z, self.z+4
    do
      setBlock(UpdateQueue,px,GROUND_MONITOR_ACTIVE_LEVEL + 1,pz,E_BLOCK_WOOL,metaPrimaryColor)
    end
  end

  for py = GROUND_MONITOR_ACTIVE_LEVEL+2, GROUND_MONITOR_ACTIVE_LEVEL+3
  do
    setBlock(UpdateQueue,self.x+1,py,self.z,E_BLOCK_WOOL,metaPrimaryColor)

    -- leave empty space for the door
    -- setBlock(UpdateQueue,self.x+2,py,self.z,E_BLOCK_WOOL,metaPrimaryColor)

    setBlock(UpdateQueue,self.x,py,self.z,E_BLOCK_WOOL,metaPrimaryColor)
    setBlock(UpdateQueue,self.x+3,py,self.z,E_BLOCK_WOOL,metaPrimaryColor)

    setBlock(UpdateQueue,self.x,py,self.z+1,E_BLOCK_WOOL,metaSecondaryColor)
    setBlock(UpdateQueue,self.x+3,py,self.z+1,E_BLOCK_WOOL,metaSecondaryColor)

    setBlock(UpdateQueue,self.x,py,self.z+2,E_BLOCK_WOOL,metaPrimaryColor)
    setBlock(UpdateQueue,self.x+3,py,self.z+2,E_BLOCK_WOOL,metaPrimaryColor)

    setBlock(UpdateQueue,self.x,py,self.z+3,E_BLOCK_WOOL,metaSecondaryColor)
    setBlock(UpdateQueue,self.x+3,py,self.z+3,E_BLOCK_WOOL,metaSecondaryColor)

    setBlock(UpdateQueue,self.x,py,self.z+4,E_BLOCK_WOOL,metaPrimaryColor)
    setBlock(UpdateQueue,self.x+3,py,self.z+4,E_BLOCK_WOOL,metaPrimaryColor)

    setBlock(UpdateQueue,self.x+1,py,self.z+4,E_BLOCK_WOOL,metaPrimaryColor)
    setBlock(UpdateQueue,self.x+2,py,self.z+4,E_BLOCK_WOOL,metaPrimaryColor)
  end

  -- torch
  setBlock(UpdateQueue,self.x+1,GROUND_MONITOR_ACTIVE_LEVEL+3,self.z+3,E_BLOCK_TORCH,E_META_TORCH_ZP)

  -- start / stop lever
  setBlock(UpdateQueue,self.x+1,GROUND_MONITOR_ACTIVE_LEVEL + 3,self.z + 2,E_BLOCK_WALLSIGN,E_META_CHEST_FACING_XP)
  updateSign(UpdateQueue,self.x+1,GROUND_MONITOR_ACTIVE_LEVEL + 3,self.z + 2,"","START/STOP","---->","",2)


  if running
  then
    setBlock(UpdateQueue,self.x+1,GROUND_MONITOR_ACTIVE_LEVEL+3,self.z+1,E_BLOCK_LEVER,1)
  else
    setBlock(UpdateQueue,self.x+1,GROUND_MONITOR_ACTIVE_LEVEL+3,self.z+1,E_BLOCK_LEVER,9)
  end


  -- remove button

  setBlock(UpdateQueue,self.x+2,GROUND_MONITOR_ACTIVE_LEVEL + 3,self.z + 2,E_BLOCK_WALLSIGN,E_META_CHEST_FACING_XM)
  updateSign(UpdateQueue,self.x+2,GROUND_MONITOR_ACTIVE_LEVEL + 3,self.z + 2,"","REMOVE","---->","",2)

  setBlock(UpdateQueue,self.x+2,GROUND_MONITOR_ACTIVE_LEVEL+3,self.z+3,E_BLOCK_STONE_BUTTON,E_BLOCK_BUTTON_XM)


  -- door
  -- Cuberite bug with Minecraft 1.8 apparently, doors are not displayed correctly
  -- setBlock(UpdateQueue,self.x+2,GROUND_MONITOR_ACTIVE_LEVEL+2,self.z,E_BLOCK_WOODEN_DOOR,E_META_CHEST_FACING_ZM)


  for px=self.x, self.x+3
  do
    for pz=self.z, self.z+4
    do
      setBlock(UpdateQueue,px,GROUND_MONITOR_ACTIVE_LEVEL + 4,pz,E_BLOCK_WOOL,metaPrimaryColor)
    end
  end

  setBlock(UpdateQueue,self.x+3,GROUND_MONITOR_ACTIVE_LEVEL + 2,self.z - 1,E_BLOCK_WALLSIGN,E_META_CHEST_FACING_ZM)
  updateSign(UpdateQueue,self.x+3,GROUND_MONITOR_ACTIVE_LEVEL + 2,self.z - 1,self.name,self.realname,"duitang.cn", "",2)

  -- Capacity sign
  setBlock(UpdateQueue,self.x,GROUND_MONITOR_ACTIVE_LEVEL + 2,self.z - 1,E_BLOCK_WALLSIGN,E_META_CHEST_FACING_ZM)

  -- Available sign
  setBlock(UpdateQueue,self.x+1,GROUND_MONITOR_ACTIVE_LEVEL + 2,self.z - 1,E_BLOCK_WALLSIGN,E_META_CHEST_FACING_ZM)
end

-- ActiveInstanceContainer:updateMemSign updates the mem usage
-- value displayed on ActiveInstanceContainer's sign
function ActiveInstanceContainer:updateCapacitySign(s)
  updateSign(UpdateQueue,self.x,GROUND_MONITOR_ACTIVE_LEVEL + 2,self.z - 1,"Capacity","",s,"",2)
end

-- ActiveInstanceContainer:updateCPUSign updates the mem usage
-- value displayed on ActiveInstanceContainer's sign
function ActiveInstanceContainer:updateAvailabilitySign(s)
  updateSign(UpdateQueue,self.x+1,GROUND_MONITOR_ACTIVE_LEVEL + 2,self.z - 1,"Availability","",s,"",2)
end

-- ActiveInstanceContainer:addGround creates ground blocks
-- necessary to display the container
function ActiveInstanceContainer:addGround()
  local y = GROUND_MONITOR_ACTIVE_LEVEL
  local max_x = GROUND_MAX_X

  if GROUND_MIN_X > self.x - 2
  then
    max_x = GROUND_MIN_X
    GROUND_MIN_X = self.x - 2
    min_x = GROUND_MIN_X
  end

  local min_x = GROUND_MIN_X
  for x= min_x, max_x
  do
    for z=GROUND_MIN_Z,GROUND_MAX_Z
    do
      setBlock(UpdateQueue,x,y,z,E_BLOCK_GRASS,0)
      for sky=y+1,y+6
      do
        setBlock(UpdateQueue,x,sky,z,E_BLOCK_AIR,0)
      end
    end
  end
end
