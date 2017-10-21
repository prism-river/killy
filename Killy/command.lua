function KillyCommand(Split, Player)
  if table.getn(Split) > 0
  then
    LOG("Split[1]: " .. Split[1])

    if Split[1] == "/killy"
    then
      table.remove(Split,1)
      LOG(Split)
      SendTCPMessage("query","no-args",Split,0)
    end
  end

  return true
end
