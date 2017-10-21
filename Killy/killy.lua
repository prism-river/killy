----------------------------------------
-- GLOBALS
----------------------------------------

-- queue containing the updates that need to be applied to the minecraft world
UpdateQueue = nil
-- array of container objects
Containers = {}
--
SignsToUpdate = {}
-- as a lua array cannot contain nil values, we store references to this object
-- in the "Containers" array to indicate that there is no container at an index
EmptyContainerSpace = {}

localPlayer = nil

----------------------------------------
-- FUNCTIONS
----------------------------------------

-- Tick is triggered by cPluginManager.HOOK_TICK
function Tick(TimeDelta)
  UpdateQueue:update(MAX_BLOCK_UPDATE_PER_TICK)
end

-- Plugin initialization
function Initialize(Plugin)
  Plugin:SetName("Killy")
  Plugin:SetVersion(1)

  UpdateQueue = NewUpdateQueue()

  -- Hooks

  cPluginManager:AddHook(cPluginManager.HOOK_PLAYER_JOINED, PlayerJoined);
  cPluginManager:AddHook(cPluginManager.HOOK_PLAYER_USING_BLOCK, PlayerUsingBlock);
  cPluginManager:AddHook(cPluginManager.HOOK_PLAYER_FOOD_LEVEL_CHANGE, OnPlayerFoodLevelChange);
  cPluginManager:AddHook(cPluginManager.HOOK_TAKE_DAMAGE, OnTakeDamage);
  cPluginManager:AddHook(cPluginManager.HOOK_WEATHER_CHANGING, OnWeatherChanging);
  cPluginManager:AddHook(cPluginManager.HOOK_SERVER_PING, OnServerPing);
  cPluginManager:AddHook(cPluginManager.HOOK_TICK, Tick);

  -- Command Bindings

  -- TODO
  cPluginManager.BindCommand("/killy", "*", KillyCommand, " - docker CLI commands")

  -- make all players admin
  cRankManager:SetDefaultRank("Admin")

  cNetwork:Connect("127.0.0.1",25566,TCP_CLIENT)

  LOG("Initialised " .. Plugin:GetName() .. " v." .. Plugin:GetVersion())

  return true
end

-- getStartStopLeverContainer returns the container
-- id that corresponds to lever at x,y coordinates
function getStartStopLeverContainer(x, z)
  for i=1, table.getn(Containers)
  do
    if Containers[i] ~= EmptyContainerSpace and x == Containers[i].x + 1 and z == Containers[i].z + 1
    then
      return Containers[i].id
    end
  end
  return ""
end

-- getRemoveButtonContainer returns the container
-- id and state for the button at x,y coordinates
function getRemoveButtonContainer(x, z)
  for i=1, table.getn(Containers)
  do
    if Containers[i] ~= EmptyContainerSpace and x == Containers[i].x + 2 and z == Containers[i].z + 3
    then
      return Containers[i].id, Containers[i].running
    end
  end
  return "", true
end

--
function PlayerJoined(Player)
  -- enable flying
  Player:SetCanFly(true)
  LOG("player joined")
  localPlayer = Player
  -- updateTableRecordContainer(1,"?", "??")
  -- updateTableRecordContainer(2,"!", "!!")
  -- updateActiveInstanceContainer(1,"??",true)
  -- updateActiveInstanceContainer(2,"!!",true)
end

--
function PlayerUsingBlock(Player, BlockX, BlockY, BlockZ, BlockFace, CursorX, CursorY, CursorZ, BlockType, BlockMeta)
  LOG("Using block: " .. tostring(BlockX) .. "," .. tostring(BlockY) .. "," .. tostring(BlockZ) .. " - " .. tostring(BlockType) .. " - " .. tostring(BlockMeta))
end

function OnPlayerFoodLevelChange(Player, NewFoodLevel)
  -- Don't allow the player to get hungry
  return true, Player, NewFoodLevel
end

function OnTakeDamage(Receiver, TDI)
  -- Don't allow the player to take falling or explosion damage
  if Receiver:GetClass() == 'cPlayer'
  then
    if TDI.DamageType == dtFall or TDI.DamageType == dtExplosion then
      return true, Receiver, TDI
    end
  end
  return false, Receiver, TDI
end

function OnServerPing(ClientHandle, ServerDescription, OnlinePlayers, MaxPlayers, Favicon)
  -- Change Server Description
  local serverDescription = "A Docker client for Minecraft"
  -- Change favicon
  if cFile:IsFile("/srv/logo.png") then
    local FaviconData = cFile:ReadWholeFile("/srv/logo.png")
    if (FaviconData ~= "") and (FaviconData ~= nil) then
      Favicon = Base64Encode(FaviconData)
    end
  end
  return false, serverDescription, OnlinePlayers, MaxPlayers, Favicon
end

-- Make it sunny all the time!
function OnWeatherChanging(World, Weather)
  return true, wSunny
end

