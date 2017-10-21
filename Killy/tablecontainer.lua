-- Container object is the representation of a Docker
-- container in the Minecraft world

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

-- TableRecordContainer:display displays all Container's blocks
-- Blocks will be blue if the container is running,
-- orange otherwise.
function TableRecordContainer:display()

  local metaPrimaryColor = E_META_WOOL_LIGHTBLUE
  local metaSecondaryColor = E_META_WOOL_BLUE

  if self.id == 0
  then
    metaPrimaryColor = E_META_WOOL_ORANGE
    metaSecondaryColor = E_META_WOOL_RED
  end

  if self.name == "query"
  then 
    metaPrimaryColor = E_META_WOOL_GREEN
  end

  self.displayed = true

  local counter = 1
  for i=1, table.getn(self.content), 4
  do
    -- add a block and add a sign to the blocks
    setBlock(UpdateQueue,self.x+1,GROUND_TABLE_LEVEL+counter,self.z + 1,E_BLOCK_WOOL,metaPrimaryColor) 
    setBlock(UpdateQueue,self.x,GROUND_TABLE_LEVEL+counter,self.z + 1,E_BLOCK_WALLSIGN,E_META_CHEST_FACING_XM)
    local first = self.content[i]
    local second = self.content[i+1]
    local third = self.content[i+2]
    local fourth = self.content[i+3]
    if first == nil
    then
      first = ""
    end
    if second == nil
    then
      second = ""
    end
    if third == nil
    then
      third = ""
    end
    if fourth == nil
    then
      fourth = ""
    end
    updateSign(UpdateQueue,self.x,GROUND_TABLE_LEVEL+counter,self.z + 1,first, second, third, fourth,0)
    counter = counter + 1
  end
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
      setBlock(UpdateQueue,x,y,z,E_BLOCK_WOOL,E_META_WOOL_WHITE)
      for sky=y+1,y+6
      do
        setBlock(UpdateQueue,x,sky,z,E_BLOCK_AIR,0)
      end
    end
  end
end
