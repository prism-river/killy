-- Container object is the representation of a Docker
-- container in the Minecraft world

-- NewTableRecordContainer returns a Container object,
-- representation of a container in
-- the Minecraft world
function NewTableRecordContainer()
  c = {
    displayed = false,
    x = 0,
    y = 0,
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
    y = 0,
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

-- TableRecordContainer:display displays all Container's blocks
-- Blocks will be blue if the container is running,
-- orange otherwise.
function TableRecordContainer:display()

  local metaPrimaryColor = E_META_WOOL_LIGHTBLUE
  local metaSecondaryColor = E_META_WOOL_BLUE
  local level = 0

  if self.name == "column"
  then
    LOG("!!!!!!!!!!!!!!!!!!!!")
    LOG(self.x)
    LOG(self.z)
    TABLE_SIGNAL_OFFSET = 0
    metaPrimaryColor = E_META_CONCRETE_GRAY
    metaSecondaryColor = E_META_CONCRETE_GRAY
    for i=1, table.getn(self.content)
    do
      -- add a block and add a sign to the blocks
      setBlock(UpdateQueue,self.x+i,GROUND_TABLE_LEVEL + 1,self.z,E_BLOCK_WOOL,metaPrimaryColor) 
      setBlock(UpdateQueue,self.x+i,GROUND_TABLE_LEVEL + 1,self.z - 1,E_BLOCK_WALLSIGN,E_META_CHEST_FACING_ZM)
      updateSign(UpdateQueue,self.x+i,GROUND_TABLE_LEVEL + 1,self.z - 1,"", self.content[i], "", "",0)
    end
  elseif self.name == "query"
  then 
    TABLE_SIGNAL_OFFSET = 0
    local level = tonumber(self.id)
    metaPrimaryColor = E_META_WOOL_GREEN
    for i=1, table.getn(self.content)
    do
      -- add a block and add a sign to the blocks
      setBlock(UpdateQueue,self.x+i,GROUND_TABLE_LEVEL + 1 + level,self.z,E_BLOCK_WOOL,metaPrimaryColor) 
      setBlock(UpdateQueue,self.x+i,GROUND_TABLE_LEVEL + 1 + level,self.z - 1,E_BLOCK_WALLSIGN,E_META_CHEST_FACING_ZM)
      updateSign(UpdateQueue,self.x+i,GROUND_TABLE_LEVEL + 1 + level,self.z - 1,"", self.content[i], "", "",0)
    end
  else
    LOG("??????")
    LOG(self.x)
    LOG(self.z)
    TABLE_SIGNAL_OFFSET = 0
    local level = tonumber(self.name)
    for i=1, table.getn(self.content)
    do
      -- add a block and add a sign to the blocks
      setBlock(UpdateQueue,self.x+i,GROUND_TABLE_LEVEL + 1 + level,self.z,E_BLOCK_WOOL,metaPrimaryColor) 
      setBlock(UpdateQueue,self.x+i,GROUND_TABLE_LEVEL + 1 + level,self.z - 1,E_BLOCK_WALLSIGN,E_META_CHEST_FACING_ZM)
      updateSign(UpdateQueue,self.x+i,GROUND_TABLE_LEVEL + 1 + level,self.z - 1,"", self.content[i], "", "",0)
    end
  end

  self.displayed = true
end

-- TableRecordContainer:addGround creates ground blocks
-- necessary to display the container
function TableRecordContainer:addGround()
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
