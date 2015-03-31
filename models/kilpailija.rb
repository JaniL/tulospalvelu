require './lib/database.rb'
require './lib/base_model.rb'

class Kilpailija < BaseModel

  @@sarakkeet = ['id','nimi','kansallisuus','sukupuoli','syntynyt']

  def initialize(rivi=nil)
    super rivi
  end






end