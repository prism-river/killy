-- Container object is the representation of a Docker
-- container in the Minecraft world

-- constant variables
CONTAINER_CREATED = 0
CONTAINER_RUNNING = 1
CONTAINER_STOPPED = 2

-- NewTableRecordContainer returns a Container object,
-- representation of a container in
-- the Minecraft world
function NewTableRecordContainer()
  c = {
    displayed = false,
    x = 0,
    z = 0,
    name="",
    id="",
    init=TableRecordContainer.init,
    setInfos=TableRecordContainer.setInfos,
    destroy=TableRecordContainer.destroy,
    display=TableRecordContainer.display,
    addGround=TableRecordContainer.addGround
  }
  return c
end

TableRecordContainer = {
    displayed = false, 
    x = 0,
    z = 0, 
    name="",
    id=""
}

-- TableRecordContainer:init sets Container's position
function TableRecordContainer:init(x,z)
  self.x = x
  self.z = z
  self.displayed = false
end

-- TableRecordContainer:setInfos sets Container's id, name, imageRepo,
-- image tag and running state
function TableRecordContainer:setInfos(id,name,content)
  self.id = id
  self.name = name
  self.content = content
end

-- -- TableRecordContainer:destroy removes all blocks of the
-- -- TableRecordContainer, it won't be visible on the map anymore
function TableRecordContainer:destroy(running)
  local X = self.x+2
  local Y = GROUND_TABLE_LEVEL+2
  local Z = self.z+2
  LOG("Exploding at X:" .. X .. " Y:" .. Y .. " Z:" .. Z)
  local World = cRoot:Get():GetDefaultWorld()
  World:BroadcastSoundEffect("random.explode", X, Y, Z, 1, 1)
  World:BroadcastParticleEffect("hugeexplosion",X, Y, Z, 0, 0, 0, 1, 1)

  -- if a block is removed before it's button/lever/sign, that object will drop
  -- and the player can collect it. Remove these first

  -- lever
  digBlock(UpdateQueue,self.x+1,GROUND_TABLE_LEVEL+3,self.z+1)
  -- signs
  digBlock(UpdateQueue,self.x+3,GROUND_TABLE_LEVEL+2,self.z-1)
  digBlock(UpdateQueue,self.x,GROUND_TABLE_LEVEL+2,self.z-1)
  digBlock(UpdateQueue,self.x+1,GROUND_TABLE_LEVEL+2,self.z-1)
  -- torch
  digBlock(UpdateQueue,self.x+1,GROUND_TABLE_LEVEL+3,self.z+1)
  --button
  digBlock(UpdateQueue,self.x+2,GROUND_TABLE_LEVEL+3,self.z+2)

  -- rest of the blocks
  for py = GROUND_TABLE_LEVEL+1, GROUND_TABLE_LEVEL+4
  do
    for px=self.x-1, self.x+4
    do
      for pz=self.z-1, self.z+5
      do
        digBlock(UpdateQueue,px,py,pz)
      end
    end
  end
end

-- TableRecordContainer:display displays all Container's blocks
-- Blocks will be blue if the container is running,
-- orange otherwise.
function TableRecordContainer:display()

  local metaPrimaryColor = E_META_WOOL_LIGHTBLUE
  local metaSecondaryColor = E_META_WOOL_BLUE

  self.displayed = true
  -- add a block and add a sign to the blocks
  setBlock(UpdateQueue,self.x+1,GROUND_TABLE_LEVEL + 1,self.z + 1,E_BLOCK_WOOL,metaPrimaryColor) 
  -- setBlock(UpdateQueue,self.x,GROUND_TABLE_LEVEL + 1,self.z + 2,E_BLOCK_WOOL,metaPrimaryColor)
  setBlock(UpdateQueue,self.x,GROUND_TABLE_LEVEL + 1,self.z + 1,E_BLOCK_WALLSIGN,E_META_CHEST_FACING_XM)
  updateSign(UpdateQueue,self.x,GROUND_TABLE_LEVEL + 1,self.z + 1,"","REMOVE","---->","",0)
end

-- TableRecordContainer:addGround creates ground blocks
-- necessary to display the container
function TableRecordContainer:addGround()
  local y = GROUND_TABLE_LEVEL
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
      setBlock(UpdateQueue,x,y,z,E_BLOCK_WOOL,E_META_WOOL_WHITE)
      for sky=y+1,y+6
      do
        setBlock(UpdateQueue,x,sky,z,E_BLOCK_AIR,0)
      end
    end
  end
end
