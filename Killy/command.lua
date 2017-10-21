function KillyCommand(Split, Player)
  if table.getn(Split) > 0
  then
    LOG("Split[1]: " .. Split[1])

    if Split[1] == "/killy"
    then
      table.remove(Split,1)
      SendTCPMessage("query","no-args",table.concat(Split, " "),0)
    end
  end

  return true
end
