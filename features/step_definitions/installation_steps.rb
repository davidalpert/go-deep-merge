require 'fileutils'

Given('I have installed {string} into {string} within the current directory') do |app, path|
  install_binary_and_add_to_path(app, path)
end

Given('I have installed {string} locally into the path') do |app|
  install_binary_and_add_to_path(app, 'bin')
end

def install_binary_and_add_to_path(app, path)
  exe = File.join(aruba.root_directory, path, app)
  raise "'#{exe}' not found; did you run 'make build'?" unless File.exist?(exe)

  expanded_path = expand_path(path)

  create_directory(path)
  FileUtils.cp(exe, expanded_path)

  # TODO: find a way to assert the dest path exists without listing the folder contents to stdoout
  # run_command_and_stop("ls -la #{path}", fail_on_error: true)

  unless ENV['PATH'].split(File::PATH_SEPARATOR).include?(path)
    prepend_environment_variable "PATH", expand_path(path) + File::PATH_SEPARATOR
  end
end