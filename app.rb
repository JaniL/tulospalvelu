set :root, File.dirname(__FILE__)
set :views, Proc.new { File.join(root, "views") }
set :public_folder, Proc.new { File.join(root, "static") }

get '/' do
  erb :layout
  #"Sinatra Heroku Cedar Template - The bare minimum for a sinatra app on cedar, running rack, and using bundler."
end